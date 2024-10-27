package instruction_test

import (
	"path/filepath"
	"testing"

	"github.com/buildpacks/pack/internal/buildkit/instruction"
	path_normalizer "github.com/buildpacks/pack/internal/buildkit/instruction/utils/path_normalizer"
	h "github.com/buildpacks/pack/testhelpers"
	"github.com/heroku/color"
	"github.com/moby/buildkit/client/llb"
	v1 "github.com/opencontainers/image-spec/specs-go/v1"
	"github.com/sclevine/spec"
	"github.com/sclevine/spec/report"
)

func TestHelperFunc(t *testing.T) {
	color.Disable(true)
	defer color.Disable(false)

	spec.Run(t, "phases", testHelperFunc, spec.Report(report.Terminal{}), spec.Parallel())
}

const (
	linux = "linux"
	windows = "windows"

	workdirLinux = "/home/user1/"
	workdirWindows = `c:\home\user1\`
)

var (
	scratch = llb.Scratch()

	platformLinux = v1.Platform{OS: linux}
	platformWindows = v1.Platform{OS: windows}
)

func testHelperFunc(t *testing.T, when spec.G, it spec.S) {
	when("#StateWorkDir", testStateWorkdir(t, it), spec.Parallel())

	when("#SanitizePath", testSanitizePath(t, it), spec.Parallel())

	when("#NormalizePath", testNormalizePath(t, it), spec.Parallel())

	when("#PathNormalizerForState", testPathNormalizerForState(t, it), spec.Parallel())
}

func testStateWorkdir(t *testing.T, it spec.S) func() {
	return func() {
		it("should return expected wordkir", func ()  {
			var state = scratch.Dir(workdirLinux)
			dir, err := instruction.StateWorkdir(state, platformLinux)
	
			h.AssertNilE(t, err)
			h.AssertEq(t, dir, workdirLinux)
		})
	}
}

func testSanitizePath(t *testing.T, it spec.S) func() {
	return func() {
		it("should remove drive letter for windows", func ()  {
			dir, err := instruction.SanitizePath(workdirWindows, platformWindows)

			h.AssertNilE(t, err)
			h.AssertEq(t, dir, workdirLinux)
		})

		it("should remove drive letter for linux", func ()  {
			dir, err := instruction.SanitizePath(workdirLinux, platformLinux)

			h.AssertNilE(t, err)
			h.AssertEq(t, dir, workdirLinux)
		})
	}
}

func testNormalizePath(t *testing.T, it spec.S) func() {
	return func() {
		it("should normalize path", func ()  {
			dir, err := path_normalizer.NewPathNormalizer(workdirLinux, "./path", path_normalizer.WithPlatform(platformLinux)).Normalize()

			h.AssertNilE(t, err)
			h.AssertEq(t, dir, filepath.Join(workdirLinux, "path"))
		})

		it("should navigate to parent path", func ()  {
			dir, err := path_normalizer.NewPathNormalizer(workdirLinux, "../path", path_normalizer.WithPlatform(platformLinux)).Normalize()

			h.AssertNilE(t, err)
			h.AssertEq(t, dir, filepath.Join("/home", "path"))
		})
	}
}

func testPathNormalizerForState(t *testing.T, it spec.S) func() {
	return func () {
		it("should return normalized path of current workdir of the state", func() {
			normalizer, err := instruction.PathNormalizerForState(scratch.Dir(workdirLinux), "../path", platformLinux)
			h.AssertNilE(t, err)

			dir, err := normalizer.Normalize()
			h.AssertNilE(t, err)
			h.AssertEq(t, dir, filepath.Join("/home", "path"))
		})
	}
}