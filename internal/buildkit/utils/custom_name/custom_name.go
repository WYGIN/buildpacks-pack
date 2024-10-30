package customname

import (
	"fmt"
	"math"

	"github.com/containerd/platforms"
	"github.com/moby/buildkit/frontend/dockerfile/shell"
	ocispecs "github.com/opencontainers/image-spec/specs-go/v1"
)

var _ CustomNameFactory = (*customName)(nil)

func NewCustomName(env shell.EnvGetter, ops ...CustomNameFactoryOption) CustomNameFactory {
	cn := &customName{
		env: env,
	}
	cn.WithOptions(ops...)

	return cn
}

func (cn *customName) WithOptions(ops ...CustomNameFactoryOption) error {
	for _, op := range ops {
		op(cn)
	}

	return nil
}

func (cn *customName) CustomName(cmd string) string {
	if cn.totalCmdsCount == nil || *cn.totalCmdsCount == 0 {
		return cmd
	}

	out := "["
	if cn.platform != nil {
		out += platforms.Format(*cn.platform) + 
			formatTargetPlatform(*cn.platform, platformFromEnv(cn.env)) + 
			" "
	}

	if cn.stageName != "" {
		out += cn.stageName + " "
	}

	if !cn.ignoreCMDIndexIncrement {
		(*cn.currentCmdIndex)++
	}

	out += fmt.Sprintf("%*d/%d] ", int(1+math.Log10(float64(*cn.totalCmdsCount))), cn.currentCmdIndex, cn.totalCmdsCount)
	if cn.isOnBuildCmd {
		out += "ONBUILD "
	}

	return out + cmd
}

// formatTargetPlatform formats a secondary platform string for cross compilation cases
func formatTargetPlatform(base ocispecs.Platform, target *ocispecs.Platform) string {
	if target == nil {
		return ""
	}
	if target.OS == "" {
		target.OS = base.OS
	}
	if target.Architecture == "" {
		target.Architecture = base.Architecture
	}
	p := platforms.Normalize(*target)

	if p.OS == base.OS && p.Architecture != base.Architecture {
		archVariant := p.Architecture
		if p.Variant != "" {
			archVariant += "/" + p.Variant
		}
		return "->" + archVariant
	}
	if p.OS != base.OS {
		return "->" + platforms.Format(p)
	}
	return ""
}

func platformFromEnv(env shell.EnvGetter) *ocispecs.Platform {
	var p ocispecs.Platform
	var set bool
	for _, key := range env.Keys() {
		switch key {
		case "TARGETPLATFORM":
			v, _ := env.Get(key)
			p, err := platforms.Parse(v)
			if err != nil {
				continue
			}
			return &p
		case "TARGETOS":
			p.OS, _ = env.Get(key)
			set = true
		case "TARGETARCH":
			p.Architecture, _ = env.Get(key)
			set = true
		case "TARGETVARIANT":
			p.Variant, _ = env.Get(key)
			set = true
		}
	}
	if !set {
		return nil
	}
	return &p
}
