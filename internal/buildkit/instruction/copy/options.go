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
		if err := WithUser(user)(ao); err != nil {
			return err
		}

		return WithGroup(group)(ao)
	}
}

func WithUser(user string) func(*instruction.CopyOp) error {
	return func(ao *instruction.CopyOp) error {
		ao.Chown.User = &llb.UserOpt{Name: user}
		return nil
	}
}

func WithGroup(group string) func(*instruction.CopyOp) error {
	return func(ao *instruction.CopyOp) error {
		ao.Chown.Group = &llb.UserOpt{Name: group}
		return nil
	}
}

func WithUIDGID(uid, gid int) func(*instruction.CopyOp) error {
	return func(ao *instruction.CopyOp) error {
		if err := WithUID(uid)(ao); err != nil {
			return err
		}

		return WithGID(gid)(ao)
	}
}

func WithUID(uid int) func(*instruction.CopyOp) error {
	return func(ao *instruction.CopyOp) error {
		ao.Chown.User = &llb.UserOpt{UID: uid}
		return nil
	}
}

func WithGID(gid int) func(*instruction.CopyOp) error {
	return func(ao *instruction.CopyOp) error {
		ao.Chown.Group = &llb.UserOpt{UID: gid}
		return nil
	}
}

func WithChmod(chmod uint32) func(*instruction.CopyOp) error {
	return func(ao *instruction.CopyOp) error {
		ao.Chmod = instruction.Chmod(chmod)
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

func Exclude(patterns instruction.Exclude) func(*instruction.CopyOp) error {
	return func(ao *instruction.CopyOp) error {
		ao.Exclude = patterns
		return nil
	}
}

