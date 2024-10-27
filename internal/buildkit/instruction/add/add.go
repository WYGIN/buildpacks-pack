package add

import (
	"io/fs"

	"github.com/buildpacks/pack/internal/buildkit/instruction"
	path_normalizer "github.com/buildpacks/pack/internal/buildkit/instruction/utils/path_normalizer"
	"github.com/moby/buildkit/client/llb"
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
func (a *add) ToState() (state llb.State, err error) {
	var copyOps = a.copyOps()
	if a.options.Checksum.String() != "" {
		if len(a.sources) != 1 {
			return state, errors.Errorf("checksum can't be specified for multiple sources")
		}

		if instruction.IsHttpSource(a.sources[0]) {
			return state, errors.New("checksum can't be specified for non-HTTP(S) sources")
		}
	}

	state = a.state

	normalizer, err := instruction.PathNormalizerForState(state, a.dest, a.platform)
	if err != nil {
		return state, err
	}

	dest, err := normalizer.
		WithOptions(
			path_normalizer.WithTrailingSlash(), 
			path_normalizer.WithPlatform(a.platform),
		).
		Normalize()

	if err != nil {
		return state, err
	}

	for _, src := range a.sources {
		state = state.File(
			llb.Copy(
				llb.Local(src),
				"/",
				dest,
				copyOps,
			),
		)
	}

	return state, nil
}

func (a *add) copyOps() (ops *llb.CopyInfo) {
	ops = &llb.CopyInfo{
		FollowSymlinks: true,
		CopyDirContentsOnly: true,
		CreateDestPath: true,
		AllowWildcard: true,
		AllowEmptyWildcard: true,
		AttemptUnpack: a.options.Checksum.String() != "",
		ExcludePatterns: a.options.Exclude,
		Mode: (*fs.FileMode)(&a.options.Chmod),
		ChownOpt: &llb.ChownOpt{
			User: a.options.User,
			Group: a.options.Group,
		},
	}

	return ops
}
