package expose

import "github.com/buildpacks/pack/internal/buildkit/instruction"

func Expose(ports ...instruction.Port) func(*instruction.ExposeOp) error {
	return func(eo *instruction.ExposeOp) error {
		eo.Expose = ports
		return nil
	}
}