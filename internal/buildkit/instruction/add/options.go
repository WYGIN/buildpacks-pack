package add

import (
	"github.com/buildpacks/pack/internal/buildkit/instruction"
	"github.com/moby/buildkit/client/llb"
	"github.com/opencontainers/go-digest"
)

func Link() AddOption {
	return func(ao *AddOp) error {
		ao.Link = true
		return nil
	}
}

func Exclude(patterns ...string) AddOption {
	return func(ao *AddOp) error {
		ao.Exclude = instruction.Exclude(patterns)
		return nil
	}
}

func KeepGitDir() AddOption {
	return func(ao *AddOp) error {
		ao.KeepGitDir = true
		return nil
	}
}

func WithChecksum(ref string) AddOption {
	return func(ao *AddOp) (err error) {
		checksum, err := digest.Parse(ref)
		ao.Checksum = instruction.Checksum{Digest: checksum}
		
		return err
	}
}

func WithUID(uid int) AddOption {
	return func(ao *AddOp) error {
		ao.Chown.User = &llb.UserOpt{UID: uid}
		return nil
	}
}

func WithUIDGID(uid, gid int) AddOption {
	return func(ao *AddOp) (err error) {
		err = WithUID(uid)(ao)
		ao.Chown.Group = &llb.UserOpt{UID: gid}

		return err
	}
}

func WithUser(user string) AddOption {
	return func(ao *AddOp) error {
		ao.Chown.User = &llb.UserOpt{Name: user}
		return nil
	}
}

func WithUserAndGroup(user, group string) AddOption {
	return func(ao *AddOp) (err error) {
		err = WithUser(user)(ao)
		ao.Chown.Group = &llb.UserOpt{Name: group}

		return err
	}
}

func WithChmod(perm uint32) AddOption {
	return func(ao *AddOp) error {
		ao.Chmod = instruction.Chmod(perm)
		return nil
	}
}
