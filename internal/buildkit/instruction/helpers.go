package instruction

import (
	"context"

	path_normalizer "github.com/buildpacks/pack/internal/buildkit/instruction/utils/path_normalizer"
	"github.com/moby/buildkit/client/llb"
	"github.com/moby/buildkit/util/system"
	ocispecs "github.com/opencontainers/image-spec/specs-go/v1"
	"github.com/pkg/errors"
)

func PathNormalizerForState(state llb.State, path string, platform ocispecs.Platform) (path_normalizer.PathNormalizer, error) {
	parent, err := StateWorkdir(state, platform)
	if err != nil {
		return nil, err
	}

	path, err = SanitizePath(path, platform)
	if err != nil {
		return nil, err
	}

	return path_normalizer.NewPathNormalizer(parent, path), nil
}

func SanitizePath(path string, platform ocispecs.Platform) (string, error) {
	path, err := system.CheckSystemDriveAndRemoveDriveLetter(path, platform.OS)
	if err != nil {
		return path, errors.Wrap(err, "removing drive letter")
	}

	return path, nil
}

func StateWorkdir(state llb.State, platform ocispecs.Platform) (string, error) {
	return state.GetDir(context.TODO(), llb.Platform(platform))
}
