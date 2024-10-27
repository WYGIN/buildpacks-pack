package pathnormalizer

import (
	"github.com/containerd/platforms"
	"github.com/moby/buildkit/util/system"
	ocispecs "github.com/opencontainers/image-spec/specs-go/v1"
)

func NewPathNormalizer(parent, path string, ops ...PathNormalizerOption) PathNormalizer {
	pn := &pathNormalizer{
		parent: parent,
		path: path,
		pathNormalizerOp: &pathNormalizerOp{
			platform: platforms.Normalize(platforms.DefaultSpec()),
		},
	}

	return pn.WithOptions(ops...)
}

func (pn *pathNormalizer) WithOptions(ops ...PathNormalizerOption) PathNormalizer {
	for _, op := range ops {
		op(pn.pathNormalizerOp)
	}

	return pn
}

func (pn pathNormalizer) Normalize() (string, error) {
	return normalize(pn.parent, pn.path, pn.platform, pn.keepSlash)
}

func normalize(parent, path string, platform ocispecs.Platform, keepSlash bool) (string, error) {
	if system.IsAbs(path, platform.OS) {
		return system.NormalizePath("/", path, platform.OS, keepSlash)
	}

	if path == "." || path == "" {
		path = "./"
	}

	return system.NormalizePath(parent, path, platform.OS, keepSlash)
}