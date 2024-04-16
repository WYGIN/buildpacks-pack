package buildkit

import (
	"context"

	"github.com/buildpacks/pack/internal/build"
)

type LifecycleExecutor interface {
	Execute(ctx context.Context, opts build.LifecycleOptions) error
}
