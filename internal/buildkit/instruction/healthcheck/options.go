package healthcheck

import (
	"time"

	"github.com/buildpacks/pack/internal/buildkit/instruction"
)

func WithInterval(interval time.Duration) func(*instruction.HealthcheckOp) error {
	return func(ho *instruction.HealthcheckOp) error {
		ho.Interval = interval
		return nil
	}
}

func WithTimeout(timeout time.Duration) func(*instruction.HealthcheckOp) error {
	return func(ho *instruction.HealthcheckOp) error {
		ho.Timeout = timeout
		return nil
	}
}

func WithStartPeriod(startPeriod time.Duration) func(*instruction.HealthcheckOp) error {
	return func(ho *instruction.HealthcheckOp) error {
		ho.StartPeriod = startPeriod
		return nil
	}
}

func WithStartInterval(interval time.Duration) func(*instruction.HealthcheckOp) error {
	return func(ho *instruction.HealthcheckOp) error {
		ho.StartInterval = interval
		return nil
	}
}

func WithNumOfRetries(retries int) func(*instruction.HealthcheckOp) error {
	return func(ho *instruction.HealthcheckOp) error {
		ho.Retries = retries
		return nil
	}
}
