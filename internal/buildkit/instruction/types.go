package instruction

import (
	"os"
	"time"

	"github.com/distribution/reference"
	"github.com/moby/buildkit/client/llb"
	"github.com/moby/buildkit/solver/pb"
	"github.com/opencontainers/go-digest"
	ocispecs "github.com/opencontainers/image-spec/specs-go/v1"
)

type Instruction interface {
	ToState() (llb.State, error)
}

type HistoryCommiter interface {
	Commit() error
}


type ChownChmod struct {
	Chown
	Chmod
}

type Chown llb.ChownOpt

type Chmod os.FileMode

type Link bool

type Exclude []string

type Parent bool

type Checksum struct {
	digest.Digest
}

type KeepGitDir bool

type ARG interface {
	Instruction
	WithOptions(...ArgOption) error
}

type ArgOption func(*ArgOp) error

type ArgOp struct {
	Args llb.EnvList
}

type CMD interface {
	Instruction
	WithOptions(...CmdOption) error
}

type CmdOption func(*CmdOp) error

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
	WithOptions(...CopyOption) error
	HasOption(CopyOption) bool
}

type CopyOption func(*CopyOp) error

type CopyOp struct {
	From *llb.State
	ChownChmod
	Link
	Parent bool
	Exclude

	// Source []string
	// Dest string
}

type ENTRYPOINT interface {
	Instruction
	WithOptions(...EntrypointOption) error
}

type EntrypointOption func(*EntrypointOp) error

type EntrypointOp struct {
	CmdForm
	Entrypoint []string
}

type ENV interface {
	Instruction
	WithOptions(...EnvOption) error
}

type EnvOption func(*EnvOp) error

type EnvOp struct {
	EnvList llb.EnvList
}

type EXPOSE interface {
	Instruction
	WithOptions(...ExposeOption) error
}

type ExposeOption func(*ExposeOp) error

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
	WithOptions(...FromOption) error
}

type FromOption func(*FromOp) error

type FromOp struct {
	Platform ocispecs.Platform
	Image reference.Reference
}

type HEALTHCHECK interface {
	Instruction
	WithOptions(...HealthcheckOption) error
	WithCMD(cmd CMD)
}

type HealthcheckOption func(*HealthcheckOp) error

type HealthcheckOp struct {
	Interval time.Duration
	Timeout time.Duration
	StartPeriod time.Duration
	StartInterval time.Duration
	Retries int
}

type LABEL interface {
	Instruction
	WithOptions(...LabelOption) error
}

type LabelOption func(*LabelOp) error

type LabelOp struct {
	Labels []KeyValuePair
}

type ONBUILD interface {
	Instruction
	WithOptions(...OnBuildOption) error
}

type OnBuildOption func(*OnBuildOp) error

type OnBuildOp struct {
	CMDs []Instruction
}

type RUN interface {
	Instruction
	WithOptions(...RunOption) error
}

type RunOption func(*RunOp) error

type RunOp struct {
	Mounts []llb.MountOption
	Network pb.NetMode
	Security pb.SecurityMode
}

type SHELL interface {
	Instruction
	WithOptions(...ShellOption) error
}

type ShellOption func(*ShellOp) error

type ShellOp struct {
	Executable string
	Params []string
}

type STOPSIGNAL interface {
	Instruction
	WithOptions(...StopSignalOption) error
}

type StopSignalOption func(*StopSignalOp) error

type StopSignalOp struct {
	StopSignal os.Signal
}

type USER interface {
	Instruction
	WithOptions(...UserOption) error
}

type UserOption func(*UserOp) error

type UserOp struct {
	User, Group llb.UserOpt
}

type VOLUME interface {
	Instruction
	WithOptions(...VolumeOption) error
}

type VolumeOption func(*VolumeOp) error

type VolumeOp struct {
	Volumes []string
}

type WORKDIR interface {
	Instruction
	WithOptions(...WorkDirOption) error
}

type WorkDirOption func(*WorkDirOp) error

type WorkDirOp struct {
	WorkDir string
}

type KeyValuePair struct {
	Key, Value string
}
