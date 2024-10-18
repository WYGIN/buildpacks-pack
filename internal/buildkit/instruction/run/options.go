package run

import (
	"github.com/buildpacks/pack/internal/buildkit/instruction"
	"github.com/moby/buildkit/client/llb"
	"github.com/moby/buildkit/solver/pb"
)

func WithMounts(mounts ...llb.MountOption) func(*instruction.RunOp) error {
	return func(ro *instruction.RunOp) error {
		ro.Mounts = mounts
		return nil
	}
}

func WithSecurity(security pb.SecurityMode) func(*instruction.RunOp) error {
	return func(ro *instruction.RunOp) error {
		ro.Security = security
		return nil
	}
}

func WithNetwork(network pb.NetMode) func(*instruction.RunOp) error {
	return func(ro *instruction.RunOp) error {
		ro.Network = network
		return nil
	}
}
