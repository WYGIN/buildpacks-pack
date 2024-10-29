package add

import (
	"io/fs"

	"github.com/buildpacks/pack/internal/buildkit/instruction"
	path_normalizer "github.com/buildpacks/pack/internal/buildkit/instruction/utils/path_normalizer"
	"github.com/moby/buildkit/client/llb"
	"github.com/moby/buildkit/solver/pb"
	v1 "github.com/opencontainers/image-spec/specs-go/v1"
	"github.com/pkg/errors"
)

var _ ADD = (*add)(nil)

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
func (a add) ToState() (llb.State, error) {
	if err := validateAddCommand(a); err != nil {
		return a.state, err
	}

	dest, err := normalizedPath(a.state, a.dest, a.platform)
	if err != nil {
		return a.state, err
	}

	a.state = copySourcesToState(a, dest)

	err = NewAddHistoryCommitter(&a).Commit()
	return a.state, err
}

func copySourcesToState(a add, dest string) llb.State {
	var fa *llb.FileAction
	for _, src := range a.sources {
		var srcAndDest = SourceAndDest{src: src, dest: dest}
		copier := NewSourceCopier(a, srcAndDest)
		copier.Copy(func(state llb.State, src, dest string, ops ...llb.CopyOption) {
			if fa == nil {
				fa = llb.Copy(state, src, dest, ops...)
			} else {
				fa = fa.Copy(state, src, dest, ops...)
			}
		})
	}

	fileOpt := []llb.ConstraintsOpt{
		// llb.WithCustomName(pgName),
		// location(cfg.opt.sourceMap, cfg.location),
	}
	// if d.ignoreCache {
	// 	fileOpt = append(fileOpt, llb.IgnoreCache)
	// }

	caps := a.client.BuildOpts().Caps
	if (&caps).Supports(pb.CapMergeOp) == nil && bool(a.options.Link) && a.options.Chmod.IsNil() {
		a.state = linkState(&a, fa, fileOpt)
	} else {
		a.state = a.state.File(fa, fileOpt...)
	}

	return a.state
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

func linkState(add *add, fa *llb.FileAction, fileOpt []llb.ConstraintsOpt) llb.State {
	// pgID := identity.NewID()
	// d.cmdIndex-- // prefixCommand increases it
	// pgName := prefixCommand(d, name, d.prefixPlatform, &platform, env)

	copyOpts := []llb.ConstraintsOpt{
		llb.Platform(add.platform),
	}
	copyOpts = append(copyOpts, fileOpt...)
	// copyOpts = append(copyOpts, llb.ProgressGroup(pgID, pgName, true))

	mergeOpts := append([]llb.ConstraintsOpt{}, fileOpt...)
	// d.cmdIndex--
	// mergeOpts = append(mergeOpts, llb.ProgressGroup(pgID, pgName, false), llb.WithCustomName(prefixCommand(d, "LINK "+name, d.prefixPlatform, &platform, env)))

	add.state = add.state.WithOutput(llb.Merge([]llb.State{add.state, llb.Scratch().File(fa, copyOpts...)}, mergeOpts...).Output())

	return add.state
}
