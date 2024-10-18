package copy

import (
	"github.com/buildpacks/pack/internal/buildkit/instruction"
	"github.com/moby/buildkit/client/llb"
)

func From(state *llb.State) func(*instruction.CopyOp) error {
	return func(co *instruction.CopyOp) error {
		co.From = state
		return nil
	}
}

func WithUserAndGroup(user, group string) func(*instruction.CopyOp) error {
	return func(ao *instruction.CopyOp) error {
		ao.Chown = llb.ChownOpt{
			User: &llb.UserOpt{Name: user},
			Group: &llb.UserOpt{Name: group},
		}

		return nil
	}
}

func WithUIDGID(uid, gid int) func(*instruction.CopyOp) error {
	return func(ao *instruction.CopyOp) error {
		ao.Chown = llb.ChownOpt{
			User: &llb.UserOpt{UID: uid},
			Group: &llb.UserOpt{UID: gid},
		}

		return nil
	}
}

func WithChmod(chmod ...uint16) func(*instruction.CopyOp) error {
	return func(ao *instruction.CopyOp) error {
		ao.Chmod = append(ao.Chmod, chmod...)
		return nil
	}
}

func Link() func(*instruction.CopyOp) error {
	return func(ao *instruction.CopyOp) error {
		ao.Link = true
		return nil
	}
}

func Parent()  func(*instruction.CopyOp) error {
	return func(ao *instruction.CopyOp) error {
		ao.Parent = true
		return nil
	}
}

func Exclude(patterns ...string) func(*instruction.CopyOp) error {
	return func(ao *instruction.CopyOp) error {
		ao.Exclude = patterns
		return nil
	}
}

