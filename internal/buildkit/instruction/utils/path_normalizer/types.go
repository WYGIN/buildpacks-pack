package pathnormalizer

import ocispecs "github.com/opencontainers/image-spec/specs-go/v1"

type PathNormalizer interface {
	WithOptions(...PathNormalizerOption) PathNormalizer
	Normalize() (string, error)
}

type PathNormalizerOption func(*pathNormalizerOp)

type pathNormalizerOp struct {
	keepSlash bool
	platform ocispecs.Platform
}

var _ PathNormalizer = (*pathNormalizer)(nil)

type pathNormalizer struct {
	parent, path string
	*pathNormalizerOp
}