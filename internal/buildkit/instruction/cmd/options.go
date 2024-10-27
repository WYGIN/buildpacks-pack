package cmd

import "github.com/buildpacks/pack/internal/buildkit/instruction"

func WithCmd(cmd []string) func(*instruction.CmdOp) error {
	return func(co *instruction.CmdOp) error {
		co.CmdForm = instruction.Exec
		co.Cmd = cmd
		return nil
	}
}

func WithShellCmd(cmd []string) func(*instruction.CmdOp) error {
	return func(co *instruction.CmdOp) error {
		co.CmdForm = instruction.Shell
		co.Cmd = cmd
		return nil
	}
}