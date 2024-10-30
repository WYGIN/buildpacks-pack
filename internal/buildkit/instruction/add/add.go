package add

import (
	"context"
	"fmt"
	"io/fs"

	"github.com/buildpacks/pack/internal/buildkit/instruction"
	path_normalizer "github.com/buildpacks/pack/internal/buildkit/instruction/utils/path_normalizer"
	customname "github.com/buildpacks/pack/internal/buildkit/utils/custom_name"
	"github.com/containerd/platforms"
	"github.com/moby/buildkit/client/llb"
	"github.com/moby/buildkit/identity"
	"github.com/moby/buildkit/solver/pb"
	v1 "github.com/opencontainers/image-spec/specs-go/v1"
	"github.com/pkg/errors"
)

var _ ADD = (*add)(nil)

func (a *add) init() {
	if a.platform.OS == "" {
		a.platform = platforms.DefaultSpec()
	}

	if a.options == nil {
		a.options = new(AddOp)
	}

	if a.currentCmdIndex == nil {
		a.currentCmdIndex = new(uint)
	}

	if a.totalCmdCount == nil {
		a.currentCmdIndex = new(uint)
	}

	if *a.currentCmdIndex > *a.totalCmdCount {
		panic("currentCommandIndex cannot be greater than totalCommandsCount")
	}

	if isAddCommand(*a) {
		a.cmd = "ADD"
	} else {
		a.cmd = "COPY"
	}

	if len(a.sources) == 0 {
		panic(fmt.Sprintf("no sources provided for command %s", a.cmd))
	}
}

func Add(state llb.State, path SourcesAndDest, platform v1.Platform) *add {
	return &add{
		state: state,
		sources: path.sources,
		options: new(AddOp),
		dest: path.dest,
		platform: platform,
	}
}

// WithOptions implements ADD.
func (a *add) WithOptions(ops ...AddOption) error {
	for _, op := range ops {
		if err := op(a.options); err != nil {
			return errors.Wrap(err, "add option error")
		}
	}

	return nil
}

// ToState implements ADD.
func (a *add) ToState() (llb.State, error) {
	if err := validateAddCommand(*a); err != nil {
		return a.state, err
	}

	dest, err := normalizedPath(a.state, a.dest, a.platform)
	if err != nil {
		return a.state, err
	}

	a.state, err = copySourcesToState(a, dest)
	if err != nil {
		return a.state, err
	}

	err = NewAddHistoryCommitter(a).Commit()
	return a.state, err
}

func copySourcesToState(a *add, dest string) (llb.State, error) {
	var fa *llb.FileAction
	for _, src := range a.sources {
		var srcAndDest = SourceAndDest{src: src, dest: dest}
		copier := NewSourceCopier(*a, srcAndDest)
		copier.Copy(func(state llb.State, src, dest string, ops ...llb.CopyOption) {
			if fa == nil {
				fa = llb.Copy(state, src, dest, ops...)
			} else {
				fa = fa.Copy(state, src, dest, ops...)
			}
		})
	}

	cn, err := buildCustomNameImpl(a, true)
	if err != nil {
		return a.state, err
	}

	fileOpt := []llb.ConstraintsOpt{
		llb.WithCustomName(cn.CustomName(a.cmd)),
		// location(cfg.opt.sourceMap, cfg.location),
	}
	if a.options.ignoreCache {
		fileOpt = append(fileOpt, llb.IgnoreCache)
	}

	caps := a.client.BuildOpts().Caps
	if (&caps).Supports(pb.CapMergeOp) == nil && bool(a.options.Link) && a.options.Chmod.IsNil() {
		a.state, err = linkState(a, fa, fileOpt)
	} else {
		a.state = a.state.File(fa, fileOpt...)
	}

	return a.state, err
}

func (a add) copyOps(src string) (ops *llb.CopyInfo) {
	a.sources = []string{src}
	ops = &llb.CopyInfo{
		FollowSymlinks: true,
		CopyDirContentsOnly: true,
		CreateDestPath: true,
		AllowWildcard: true,
		AllowEmptyWildcard: true,
		AttemptUnpack: isAddCommand(a),
		ExcludePatterns: a.options.Exclude,
		Mode: (*fs.FileMode)(&a.options.Chmod),
		ChownOpt: (*llb.ChownOpt)(&a.options.Chown),
	}

	return ops
}

func validateAddCommand(a add) error {
	if a.options.Checksum.String() == "" {
		return nil
	}

	if len(a.sources) != 1 {
		return errors.Errorf("checksum can't be specified for multiple sources")
	}

	if instruction.IsHttpSource(a.sources[0]) {
		return nil
	}

	return errors.New("checksum can't be specified for non-HTTP(S) sources")
}

func normalizedPath(state llb.State, dest string, platform v1.Platform) (string,error) {
	normalizer, err := instruction.PathNormalizerForState(state, dest, platform)
	if err != nil {
		return dest, err
	}

	return normalizer.
		WithOptions(
			path_normalizer.WithTrailingSlash(), 
			path_normalizer.WithPlatform(platform),
		).
		Normalize()
}

func linkState(add *add, fa *llb.FileAction, fileOpt []llb.ConstraintsOpt) (llb.State, error) {
	pgID := identity.NewID()

	cn, err := buildCustomNameImpl(add, false)
	if err != nil {
		return add.state, err
	}

	pgName := cn.CustomName(add.cmd)

	copyOpts := []llb.ConstraintsOpt{
		llb.Platform(add.platform),
	}
	copyOpts = append(copyOpts, fileOpt...)
	copyOpts = append(copyOpts, llb.ProgressGroup(pgID, pgName, true))

	mergeOpts := append([]llb.ConstraintsOpt{}, fileOpt...)
	// d.cmdIndex--
	mergeOpts = append(mergeOpts, llb.ProgressGroup(pgID, pgName, false), llb.WithCustomName(cn.CustomName("LINK " + add.cmd)))

	add.state = add.state.WithOutput(llb.Merge([]llb.State{add.state, llb.Scratch().File(fa, copyOpts...)}, mergeOpts...).Output())

	return add.state, nil
}

func buildCustomNameImpl(add *add, cmdIndexIncrement bool) (customname.CustomNameFactory, error) {
	envList, err := add.state.Env(context.TODO())
	if err != nil {
		return nil, err
	}

	customnameOps := []customname.CustomNameFactoryOption{
		customname.CurrentCommandIndex(add.currentCmdIndex),
		customname.TotalCmdsCount(add.totalCmdCount),
		customname.WithPlatform(add.platform),
	}

	if !cmdIndexIncrement {
		customnameOps = append(customnameOps, customname.IgnoreCmdIndexIncrement)
	}

	if add.options.onBuildCmd {
		customnameOps = append(customnameOps, customname.OnBuildCMD)
	}

	cn := customname.NewCustomName(
		envList,
		customnameOps...
	) // prefixCommand(a, name, a.prefixPlatform, &platform, env)

	return cn, nil
}
