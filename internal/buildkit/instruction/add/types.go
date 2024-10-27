package add

import (
	"fmt"
	"strings"

	"github.com/buildpacks/pack/internal/buildkit/instruction"
	"github.com/moby/buildkit/client/llb"
	ocispecs "github.com/opencontainers/image-spec/specs-go/v1"
)

type ADD interface {
	instruction.Instruction
	WithOptions(...AddOption) error
}

type AddOption func(*AddOp) error

type AddOp struct {
	instruction.KeepGitDir
	instruction.Checksum
	instruction.Chown
	instruction.Chmod
	instruction.Exclude
	instruction.Link
}

type add struct {
	state llb.State
	options *AddOp
	sources []string
	dest string
	platform ocispecs.Platform
}

type SourcesAndDest struct {
	sources []string
	dest string
}

type AddHistroryCommiter interface {
	instruction.HistoryCommiter
}

type addHistroryCommiter struct {
	frontend string
	add ADD
	AddCMDStringer
}

type AddCMDStringer interface {
	fmt.Stringer
}

type addCMDStringer struct {
	builder *strings.Builder
	add
}
