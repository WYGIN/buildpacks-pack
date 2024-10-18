package add

import (
	"github.com/buildpacks/pack/internal/buildkit/instruction"
	"github.com/moby/buildkit/client/llb"
	digest "github.com/opencontainers/go-digest"
)

func KeepGitDir() func(* instruction.AddOp) error {
	return func(ao *instruction.AddOp) error {
		ao.KeepGitDir = true
		return nil
	}
}

func WithChecksum(checksum string) func(*instruction.AddOp) error {
	return func(ao *instruction.AddOp) (err error) {
		ao.Checksum, err = digest.Parse(checksum)
		return err
	}
}

func WithUserAndGroup(user, group string) func(*instruction.AddOp) error {
	return func(ao *instruction.AddOp) error {
		ao.Chown = llb.ChownOpt{
			User: &llb.UserOpt{Name: user},
			Group: &llb.UserOpt{Name: group},
		}

		return nil
	}
}

func WithUIDGID(uid, gid int) func(*instruction.AddOp) error {
	return func(ao *instruction.AddOp) error {
		ao.Chown = llb.ChownOpt{
			User: &llb.UserOpt{UID: uid},
			Group: &llb.UserOpt{UID: gid},
		}

		return nil
	}
}

func WithChmod(chmod ...uint16) func(*instruction.AddOp) error {
	return func(ao *instruction.AddOp) error {
		ao.Chmod = append(ao.Chmod, chmod...)
		return nil
	}
}

func Link() func(*instruction.AddOp) error {
	return func(ao *instruction.AddOp) error {
		ao.Link = true
		return nil
	}
}

func Exclude(patterns ...string) func(*instruction.AddOp) error {
	return func(ao *instruction.AddOp) error {
		ao.Exclude = patterns
		return nil
	}
}
