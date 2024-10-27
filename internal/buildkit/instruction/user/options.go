package user

import (
	"github.com/buildpacks/pack/internal/buildkit/instruction"
	"github.com/moby/buildkit/client/llb"
)

func WithUserAndGroup(user, group string) func(*instruction.UserOp) error {
	return func(ao *instruction.UserOp) error {
		if err := WithUser(user)(ao); err != nil {
			return err
		}

		return WithGroup(group)(ao)
	}
}

func WithUser(user string) func(*instruction.UserOp) error {
	return func(ao *instruction.UserOp) error {
		ao.User = llb.UserOpt{Name: user}
		return nil
	}
}

func WithGroup(group string) func(*instruction.UserOp) error {
	return func(ao *instruction.UserOp) error {
		ao.Group = llb.UserOpt{Name: group}
		return nil
	}
}

func WithUIDGID(uid, gid int) func(*instruction.UserOp) error {
	return func(ao *instruction.UserOp) error {
		if err := WithUID(uid)(ao); err != nil {
			return err
		}

		return WithGID(gid)(ao)
	}
}

func WithUID(uid int) func(*instruction.UserOp) error {
	return func(ao *instruction.UserOp) error {
		ao.User = llb.UserOpt{UID: uid}
		return nil
	}
}

func WithGID(gid int) func(*instruction.UserOp) error {
	return func(ao *instruction.UserOp) error {
		ao.Group = llb.UserOpt{UID: gid}
		return nil
	}
}
