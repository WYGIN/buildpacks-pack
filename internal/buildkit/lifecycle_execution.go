package buildkit

import (
	"github.com/google/go-containerregistry/pkg/name"

	"github.com/buildpacks/lifecycle/api"
	"github.com/buildpacks/pack/internal/build"
	"github.com/buildpacks/pack/internal/paths"
	"github.com/buildpacks/pack/pkg/dist"
	"github.com/buildpacks/pack/pkg/logging"
)

func NewLifecycleExecution(logger logging.Logger, state State, targets []dist.Target, opts build.LifecycleOptions) (exec *LifecycleExecution, err error) {
	var supportedPlatformAPIs = append(opts.Builder.LifecycleDescriptor().APIs.Platform.Deprecated, opts.Builder.LifecycleDescriptor().APIs.Platform.Supported...)
	latestSupportedPlatformAPI, err := build.FindLatestSupported(supportedPlatformAPIs, opts.LifecycleApis)
	exec = &LifecycleExecution{
		logger:       logger,
		layersVolume: paths.FilterReservedNames("pack-layers-" + randString(10)),
		appVolume:    paths.FilterReservedNames("pack-app-" + randString(10)),
		platformAPI:  latestSupportedPlatformAPI,
		opts:         opts,
		targets:      targets,
		state:        state,
	}
	if opts.Interactive {
		exec.logger = opts.Termui
	}
	return exec, err
}

func (l *LifecycleExecution) Builder() build.Builder {
	return l.opts.Builder
}

func (l *LifecycleExecution) AppPath() string {
	return l.opts.AppPath
}

func (l LifecycleExecution) Workspace() string {
	return l.opts.Workspace
}

func (l *LifecycleExecution) AppVolume() string {
	return l.appVolume
}

func (l *LifecycleExecution) LayersVolume() string {
	return l.layersVolume
}

func (l *LifecycleExecution) PlatformAPI() *api.Version {
	return l.platformAPI
}

func (l *LifecycleExecution) ImageName() name.Reference {
	return l.opts.Image
}

func (l *LifecycleExecution) PrevImageName() string {
	return l.opts.PreviousImage
}
