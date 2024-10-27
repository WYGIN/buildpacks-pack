package pathnormalizer

import (
	ocispecs "github.com/opencontainers/image-spec/specs-go/v1"
)

func WithPlatform(platform ocispecs.Platform) PathNormalizerOption {
	return func(pno *pathNormalizerOp) {
		pno.platform = platform
	}
}

func WithTrailingSlash() PathNormalizerOption {
	return func(pno *pathNormalizerOp) {
		pno.keepSlash = true
	}
}