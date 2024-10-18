package shell

import "github.com/buildpacks/pack/internal/buildkit/instruction"

func WithShell(executable string, params ...string) func(*instruction.ShellOp) error {
	return func(so *instruction.ShellOp) error {
		so.Executable = executable
		so.Params = params
		return nil
	}
}