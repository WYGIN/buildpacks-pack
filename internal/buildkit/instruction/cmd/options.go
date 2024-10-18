package cmd

import "github.com/buildpacks/pack/internal/buildkit/instruction"

func WithCmd(cmd []string, form instruction.CmdForm) func(*instruction.CmdOp) error {
	return func(co *instruction.CmdOp) error {
		co.CmdForm = form
		co.Cmd = cmd
		return nil
	}
}