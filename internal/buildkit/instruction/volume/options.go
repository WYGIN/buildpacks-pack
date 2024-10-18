package volume

import "github.com/buildpacks/pack/internal/buildkit/instruction"

func WithVolumes(volumes ...string) func(*instruction.VolumeOp) error {
	return func(vo *instruction.VolumeOp) error {
		vo.Volumes = volumes
		return nil
	}
}