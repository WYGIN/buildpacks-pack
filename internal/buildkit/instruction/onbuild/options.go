package onbuild

import (
	"errors"

	"github.com/buildpacks/pack/internal/buildkit/instruction"
	"github.com/buildpacks/pack/internal/buildkit/instruction/add"
	"github.com/buildpacks/pack/internal/buildkit/instruction/copy"
)

func WithInstructions(cmds ...instruction.Instruction) func(*instruction.OnBuildOp) error {
	return func(obo *instruction.OnBuildOp) (err error) {
		for _, cmd := range cmds {
			if err := WithInstruction(cmd)(obo); err != nil {
				return err
			}
		}

		return err
	}
}

func WithInstruction(cmd instruction.Instruction) func(*instruction.OnBuildOp) error {
	return func(obo *instruction.OnBuildOp) error {
		switch cmd := cmd.(type) {
			case add.ADD,
				instruction.ARG,
				instruction.CMD,
				instruction.ENTRYPOINT,
				instruction.ENV,
				instruction.EXPOSE,
				instruction.HEALTHCHECK,
				instruction.LABEL,
				instruction.RUN,
				instruction.SHELL,
				instruction.STOPSIGNAL,
				instruction.USER,
				instruction.VOLUME,
				instruction.WORKDIR:
				obo.CMDs = append(obo.CMDs, cmd)
			case instruction.COPY:
				if !cmd.HasOption(copy.From(nil)) {
					return errors.New("unsupported: [COPY] INSTRUCTION with [--from] flag")
				}
				obo.CMDs = append(obo.CMDs, cmd)
			default:
				return errors.New("unsupported: [ONBUILD] INSTRUCTION")
		}

		return nil
	}
}