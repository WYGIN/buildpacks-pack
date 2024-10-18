package expose

import "github.com/buildpacks/pack/internal/buildkit/instruction"

func Expose(ports ...instruction.Port) func(*instruction.ExposeOp) error {
	return func(eo *instruction.ExposeOp) error {
		eo.Expose = make([]instruction.Port, 0, len(ports))
		for _, port := range ports {
			eo.Expose = append(eo.Expose, instruction.Port{
				Port: port.Port, 
				Protocal: port.Protocal,
			})
		}

		return nil
	}
}