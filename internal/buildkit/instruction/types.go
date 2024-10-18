package instruction

import (
	"os"
	"time"

	"github.com/distribution/reference"
	"github.com/moby/buildkit/client/llb"
	"github.com/moby/buildkit/solver/pb"
	digest "github.com/opencontainers/go-digest"
	ocispecs "github.com/opencontainers/image-spec/specs-go/v1"
)

type Instruction interface {
	ToState() llb.State
}

type ADD interface {
	Instruction
	WithOptions(...AddOptions)
}

type AddOptions func(op *AddOp) error

type AddOp struct {
	KeepGitDir bool
	Checksum digest.Digest
	chownChmod
	link
	exclude
}

type chownChmod struct {
	Chown llb.ChownOpt
	Chmod []uint16
}

type link struct {
	Link bool
}

type exclude struct {
	Exclude []string
}

type ARG interface {
	Instruction
	WithOptions(...ArgOptions)
}

type ArgOptions func(*ArgOp) error

type ArgOp struct {
	Args llb.EnvList
}

type CMD interface {
	Instruction
	WithOptions(...CmdOptions)
}

type CmdOptions func(*CmdOp) error

type CmdOp struct {
	CmdForm
	Cmd []string
}

type CmdForm int

const (
	Exec CmdForm = iota
	Shell
)

type COPY interface {
	Instruction
	WithOptions(...CopyOptions)
	HasOption(CopyOptions) bool
}

type CopyOptions func(*CopyOp) error

type CopyOp struct {
	From *llb.State
	chownChmod
	link
	Parent bool
	exclude

	// Source []string
	// Dest string
}

type ENTRYPOINT interface {
	Instruction
	WithOptions(...EntrypointOptions)
}

type EntrypointOptions func(*EntrypointOp) error

type EntrypointOp struct {
	CmdForm
	Entrypoint []string
}

type ENV interface {
	Instruction
	WithOptions(...EnvOptions)
}

type EnvOptions func(*EnvOp) error

type EnvOp struct {
	EnvList llb.EnvList
}

type EXPOSE interface {
	Instruction
	WithOptions(...ExposeOptions)
}

type ExposeOptions func(*ExposeOp) error

type ExposeOp struct {
	Expose []Port
}

type Port struct {
	Port int
	Protocal
}

type Protocal int

const (
	Tcp Protocal = iota
	Udp
)

type FROM interface {
	Instruction
	WithOptions(...FromOptions)
}

type FromOptions func(*FromOp) error

type FromOp struct {
	Platform ocispecs.Platform
	Image reference.Reference
}

type HEALTHCHECK interface {
	Instruction
	WithOptions(...HealthcheckOptions)
	WithCMD(cmd CMD)
}

type HealthcheckOptions func(*HealthcheckOp) error

type HealthcheckOp struct {
	Interval time.Duration
	Timeout time.Duration
	StartPeriod time.Duration
	StartInterval time.Duration
	Retries int
}

type LABEL interface {
	Instruction
	WithOptions(...LabelOptions)
}

type LabelOptions func(*LabelOp) error

type LabelOp struct {
	Labels []KeyValuePair
}

type ONBUILD interface {
	Instruction
	WithOptions(...OnBuildOptions)
}

type OnBuildOptions func(*OnBuildOp) error

type OnBuildOp struct {
	CMDs []Instruction
}

type RUN interface {
	Instruction
	WithOptions(...RunOptions)
}

type RunOptions func(*RunOp) error

type RunOp struct {
	Mounts []llb.MountOption
	Network pb.NetMode
	Security pb.SecurityMode
}

type SHELL interface {
	Instruction
	WithOptions(...ShellOptions)
}

type ShellOptions func(*ShellOp) error

type ShellOp struct {
	Executable string
	Params []string
}

type STOPSIGNAL interface {
	Instruction
	WithOptions(...StopSignalOptions)
}

type StopSignalOptions func(*StopSignalOp) error

type StopSignalOp struct {
	StopSignal os.Signal
}

type USER interface{
	Instruction
	WithOptions(...UserOptions)
}

type UserOptions func(*UserOp) error

type UserOp struct {
	User, Group llb.UserOpt
}

type VOLUME interface {
	Instruction
	WithOptions(...VolumeOptions)
}

type VolumeOptions func(*VolumeOp) error

type VolumeOp struct {
	Volumes []string
}

type WORKDIR interface{
	Instruction
	WithOptions(...WorkDirOptions)
}

type WorkDirOptions func(*WorkDirOp) error

type WorkDirOp struct {
	WorkDir string
}

type KeyValuePair struct {
	Key, Value string
}

