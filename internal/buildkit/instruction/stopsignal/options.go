package stopsignal

import (
	"os"

	"github.com/buildpacks/pack/internal/buildkit/instruction"
)

func WithStopSignal(signal os.Signal) func(*instruction.StopSignalOp) error {
	return func(sso *instruction.StopSignalOp) error {
		sso.StopSignal = signal
		return nil
	}
}