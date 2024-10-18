package entrypoint

import "github.com/buildpacks/pack/internal/buildkit/instruction"

func WithEntrypoint(ep []string, form instruction.CmdForm) func(*instruction.EntrypointOp) error {
	return func(eo *instruction.EntrypointOp) error {
		eo.CmdForm = form
		eo.Entrypoint = ep
		return nil
	}
}