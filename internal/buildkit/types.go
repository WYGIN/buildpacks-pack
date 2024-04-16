package buildkit

import (
	"github.com/buildpacks/lifecycle/api"
	"github.com/moby/buildkit/client/llb"

	"github.com/buildpacks/pack/internal/build"
	"github.com/buildpacks/pack/pkg/dist"
	"github.com/buildpacks/pack/pkg/logging"
)

type LifecycleExecution struct {
	logger       logging.Logger
	state        State
	platformAPI  *api.Version
		layersVolume string
	appVolume    string
	targets      []dist.Target
	opts         build.LifecycleOptions
}

type State struct {
	llb.State
}
