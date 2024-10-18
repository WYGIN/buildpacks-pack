package user

import (
	"github.com/buildpacks/pack/internal/buildkit/instruction"
	"github.com/moby/buildkit/client/llb"
)

func WithUserAndGroup(user, group string) func(*instruction.UserOp) error {
	return func(ao *instruction.UserOp) error {
		ao.User = llb.UserOpt{Name: user}
		ao.Group = llb.UserOpt{Name: group}
		return nil
	}
}

func WithUIDGID(uid, gid int) func(*instruction.UserOp) error {
	return func(ao *instruction.UserOp) error {
		ao.User = llb.UserOpt{UID: uid}
		ao.Group = llb.UserOpt{UID: gid}
		return nil
	}
}
