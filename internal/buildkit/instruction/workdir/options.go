package workdir

import "github.com/buildpacks/pack/internal/buildkit/instruction"

func WithWorkDir(dir string) func(*instruction.WorkDirOp) error {
	return func(wdo *instruction.WorkDirOp) error {
		wdo.WorkDir = dir
		return nil
	}
}
