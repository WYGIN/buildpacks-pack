package entrypoint

import "github.com/buildpacks/pack/internal/buildkit/instruction"

func WithEntrypoint(ep ...string) func(*instruction.EntrypointOp) error {
	return func(eo *instruction.EntrypointOp) error {
		eo.CmdForm = instruction.Exec
		eo.Entrypoint = ep
		return nil
	}
}

func WithShellEntrypoint(ep ...string) func(*instruction.EntrypointOp) error {
	return func(eo *instruction.EntrypointOp) error {
		eo.CmdForm = instruction.Shell
		eo.Entrypoint = ep
		return nil
	}
}