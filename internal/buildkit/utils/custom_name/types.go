package customname

import (
	"github.com/moby/buildkit/frontend/dockerfile/shell"
	v1 "github.com/opencontainers/image-spec/specs-go/v1"
)

type CustomNameFactory interface {
	WithOptions(ops ...CustomNameFactoryOption) error
	CustomName(cmd string) string
}

type CustomNameFactoryOption func(cn *customName) error

type customName struct {
	currentCmdIndex, totalCmdsCount *uint
	platform           *v1.Platform
	env                shell.EnvGetter
	stageName          string
	isOnBuildCmd       bool
	ignoreCMDIndexIncrement  bool
}
