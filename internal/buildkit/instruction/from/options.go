package from

import (
	"github.com/buildpacks/pack/internal/buildkit/instruction"
	"github.com/distribution/reference"
	"github.com/containerd/platforms"
)

func From(ref string) func(*instruction.FromOp) error {
	return func(fo *instruction.FromOp) (err error) {
		if fo.Image, err = reference.ParseNormalizedNamed(ref); err != nil {
			return err
		}

		fo.Platform = platforms.DefaultSpec()

		return err
	}
}

func WithPlatform(ref, platform string) func(*instruction.FromOp) error {
	return func(fo *instruction.FromOp) (err error) {
		if err := From(ref)(fo); err != nil {
			return err
		}
		
		if platform != "" {
			fo.Platform, err = platforms.Parse(platform)
		}

		return err
	}
}