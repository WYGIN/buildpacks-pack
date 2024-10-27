package workdir

import (
	"path/filepath"

	"github.com/buildpacks/pack/internal/buildkit/instruction"
)

func WithWorkdir(dir string) func(*instruction.WorkDirOp) error {
	return func(wdo *instruction.WorkDirOp) error {
		wdo.WorkDir = dir
		return nil
	}
}

func WithRelativeWorkdir(dir string) func(*instruction.WorkDirOp) error {
	return func(wdo *instruction.WorkDirOp) (err error) {
		wdo.WorkDir, err = filepath.Rel(wdo.WorkDir, dir)
		return err
	}
}
