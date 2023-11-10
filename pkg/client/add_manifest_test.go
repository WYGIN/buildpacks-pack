package client_test

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/google/go-containerregistry/pkg/name"
	v1 "github.com/google/go-containerregistry/pkg/v1"
	"github.com/heroku/color"
	"github.com/sclevine/spec"
	"github.com/sclevine/spec/report"
	"github.com/spf13/cobra"

	"github.com/buildpacks/imgutil/fakes"
	"github.com/buildpacks/pack/internal/commands"
	"github.com/buildpacks/pack/internal/commands/testmocks"
	"github.com/buildpacks/pack/pkg/logging"
	h "github.com/buildpacks/pack/testhelpers"
)

func TestManifestAddCommand(t *testing.T) {
	color.Disable(true)
	defer color.Disable(false)

	spec.Run(t, "Commands", testManifestAddCommand, spec.Random(), spec.Report(report.Terminal{}))
}

func testManifestAddCommand(t *testing.T, when spec.G, it spec.S) {
	var (
		command        *cobra.Command
		logger         *logging.LogWithWriters
		outBuf         bytes.Buffer
		mockController *gomock.Controller
		mockClient     *testmocks.MockPackClient
		tmpDir			string
		manifestListName 	string
	)

	it.Before(func() {
		outBuf = bytes.Buffer{}
		logger = logging.NewLogWithWriters(&outBuf, &outBuf)
		mockController = gomock.NewController(t)
		mockClient = testmocks.NewMockPackClient(mockController)
	})

	when("#AddManifest", func() {
		it.Before(func() {
			tmpDir, err := os.MkdirTemp("", "manifest-test")
			h.AssertNil(t, err)
			os.Setenv("XDG_RUNTIME_DIR", tmpDir)
			logger = logging.NewLogWithWriters(&outBuf, &outBuf)

			command = commands.ManifestAdd(logger, mockClient)
			manifestCreate := commands.ManifestCreate(logger, mockClient)
			manifestListName = "cnbs/sample-package:hello-multiarch-universe"
			manifestCreate.SetArgs([]string{manifestListName})
			err = manifestCreate.Execute()
			h.AssertNil(t, err)
			ref, err := name.ParseReference(manifestListName)
			h.AssertNil(t, err)
			manifestListFile, err := os.ReadFile(filepath.Join(tmpDir,"manifests", ref.Name(), ref.Identifier()))
			h.AssertNil(t, err)
			h.AssertNotNil(t, manifestListFile)
		})
		it.After(func() {
			err := os.RemoveAll(os.Getenv("XDG_RUNTIME_DIR"))
			h.AssertNil(t, err)
		})
		when("no flags specified", func() {
			it("should be able to add manifest list", func() {
				command.SetArgs([]string{manifestListName, "node:latest"})
				err := command.Execute()
				h.AssertNil(t, err)
			})
			it("should be able to add manifest", func() {
				command.SetArgs([]string{manifestListName, "node@sha256:a59381eeb372ade238dcde65dce1fb6ad48c4eda288bf4e3e50b94176ee67d60"})
				err := command.Execute()
				h.AssertNil(t, err)
			})
		})
		when("when --all flags passed", func() {
			it("should print a warning if no imageIndex is passed", func() {
				command.SetArgs([]string{manifestListName, "node@sha256:a59381eeb372ade238dcde65dce1fb6ad48c4eda288bf4e3e50b94176ee67d60", "--all"})
				err := command.Execute()
				h.AssertNil(t, err)
				h.AssertEq(t, outBuf.String(), "some warning")
			})
			it("should add all images in ImageIndex if imageIndex is passed", func() {
				var manifestList fakeBadIndexStruct
				var node = "node:18.18.2-slim"
				command.SetArgs([]string{manifestListName, node, "--all"})
				var hashes = []string{
					"sha256:ef52e84aa85baadfcdfe6f40162c368c08d4def29751ed1429abe1908316b198",
					"sha256:a02f6e55993fdb04a45017f1a9bd1876dc0a3fe89a1d88e53393752f80859e22",
					"sha256:c1e67d1a099e50d37d6aef7ee2917496a99aff6876b046613ed822bf8c72d371",
					"sha256:d1beb3473334d238c62c94ed644cfac5c7df27920579a28d15d9ed85f25e87d5",
					"sha256:581d640ed5c99190e7afcf9ffb05b030a0094a3775b35bbe05328e58877dd63a",
				}
				err := command.Execute()
				h.AssertNil(t, err)
				h.AssertEq(t, outBuf.String(), "")
				ref, err := name.ParseReference(manifestListName)
				h.AssertNil(t, err)
				file, err := os.ReadFile(filepath.Join(tmpDir,"manifests", ref.Name(), ref.Identifier()))
				h.AssertNil(t, err)
				err = json.Unmarshal(file, &manifestList)
				h.AssertNil(t, err)
				for _, hash := range hashes {
					hashRef, err := v1.NewHash(hash)
					h.AssertNil(t, err)
					img, err := manifestList.Index.Image(hashRef)
					h.AssertNil(t, err)
					h.AssertNotNil(t, img)
				}
			})
		})
		when("when --os flags passed", func() {
			it("should return an error when ImageIndex is not passed", func() {
				var node = "node@sha256:ef52e84aa85baadfcdfe6f40162c368c08d4def29751ed1429abe1908316b198"
				command.SetArgs([]string{manifestListName, node, "--os", "linux"})
				err := command.Execute()
				h.AssertNotNil(t, err)
			})
			it("should add image with given os to local list if it exists in imageIndex", func() {
				var manifestList fakeBadIndexStruct
				var node = "traefik:v3.0.0-beta4-windowsservercore-1809"
				command.SetArgs([]string{manifestListName, node, "--os", "windows"})
				err := command.Execute()
				h.AssertNil(t, err)
				ref, err := name.ParseReference(manifestListName)
				h.AssertNil(t, err)
				file, err := os.ReadFile(filepath.Join(tmpDir,"manifests", ref.Name(), ref.Identifier()))
				h.AssertNil(t, err)
				err = json.Unmarshal(file, &manifestList)
				h.AssertNil(t, err)
				var hashRef = "sha256:6bbe10dec34e310f581af3d2f4b1bca020b4b4d097063e77a5301e74af5a0196"
				img, err := manifestList.Index.Image(hashRef)
				h.AssertNil(t, err)
				h.AssertNotNil(t, img)
			})
			it("should return an error when given os doesn't exists in the ImageIndex", func() {
				var manifestList fakeBadIndexStruct
				var node = "traefik:v3.0.0-beta4-windowsservercore-1809"
				command.SetArgs([]string{manifestListName, node, "--os", "linux"})
				err := command.Execute()
				h.AssertNil(t, err)
				ref, err := name.ParseReference(manifestListName)
				h.AssertNil(t, err)
				file, err := os.ReadFile(filepath.Join(tmpDir,"manifests", ref.Name(), ref.Identifier()))
				h.AssertNil(t, err)
				err = json.Unmarshal(file, &manifestList)
				h.AssertNil(t, err)
				var hashRef = "sha256:6bbe10dec34e310f581af3d2f4b1bca020b4b4d097063e77a5301e74af5a0196"
				img, err := manifestList.Index.Image(hashRef)
				h.AssertNotNil(t, err)
				h.AssertNil(t, img)
			})
		})
		when("when --arch flags passed", func() {
			it("should not return an error when os flag is not specified", func() {
				var node = "traefik:v3.0.0-beta4"
				command.SetArgs([]string{manifestListName, node, "--arch", "amd64"})
				err := command.Execute()
				h.AssertNil(t, err)
			})
			it("should return an error when arch doesn't exists in the imageIndex", func() {
				var node = "traefik:v3.0.0-beta4"
				command.SetArgs([]string{manifestListName, node, "--arch", "abc"})
				err := command.Execute()
				h.AssertNotNil(t, err)
			})
			it("should return an error when manifest is passed instead of manifestList", func() {
				var node = "traefik:v3.0.0-beta4-windowsservercore-1809"
				command.SetArgs([]string{manifestListName, node, "--arch", "arm64"})
				err := command.Execute()
				h.AssertNotNil(t, err)
			})
			it("should include the manifest of given arch from the image index", func() {
				var manifestList fakeBadIndexStruct
				var hashes = []string{
					"sha256:777e2106170c66742ddbe77f703badb7dc94d9a5b1dc2c4a01538fad9aef56bb",
					"sha256:9909a171ac287c316f771a4f1d1384df96957ed772bc39caf6deb6e3e360316f",
				}
				var node = "alpine:3.18.4"
				command.SetArgs([]string{manifestListName, node, "--arch", "arm"})
				err := command.Execute()
				h.AssertNil(t, err)
				ref, err := name.ParseReference(manifestListName)
				h.AssertNil(t, err)
				file, err := os.ReadFile(filepath.Join(tmpDir,"manifests", ref.Name(), ref.Identifier()))
				h.AssertNil(t, err)
				err = json.Unmarshal(file, &manifestList)
				h.AssertNil(t, err)
				for _, hash := range hashes {
					img, err := manifestList.Index.Image(hash)
					h.AssertNotNil(t, err)
					h.AssertNil(t, img)
				}
			})
		})
		when("when --variant flags passed", func() {
			it("should not return an error when os flag is not specified", func() {
				var node = "alpine:3.18"
				command.SetArgs([]string{manifestListName, node, "--arch", "arm", "--variant", "v7"})
				err := command.Execute()
				h.AssertNil(t, err)
			})
			it("should not return an error when arch flag is not specified", func() {
				var node = "alpine:3.18"
				command.SetArgs([]string{manifestListName, node, "--os", "linux", "--variant", "v7"})
				err := command.Execute()
				h.AssertNil(t, err)
			})
			it("should return an error when variant doesn't exists in the imageIndex", func() {
				var node = "alpine:3.18"
				command.SetArgs([]string{manifestListName, node, "--variant", "v9"})
				err := command.Execute()
				h.AssertNotNil(t, err)
			})
			it("should return an error when manifest is passed instead of manifestList", func() {
				var node = "alpine@sha256:777e2106170c66742ddbe77f703badb7dc94d9a5b1dc2c4a01538fad9aef56bb"
				command.SetArgs([]string{manifestListName, node, "--variant", "v6"})
				err := command.Execute()
				h.AssertNotNil(t, err)
			})
			it("should include the manifest of given variant from the image index", func() {
				var manifestList fakeBadIndexStruct
				var hashes = []string{
					"sha256:777e2106170c66742ddbe77f703badb7dc94d9a5b1dc2c4a01538fad9aef56bb",
				}
				var node = "alpine:3.18.4"
				command.SetArgs([]string{manifestListName, node, "--variant", "v6"})
				err := command.Execute()
				h.AssertNil(t, err)
				ref, err := name.ParseReference(manifestListName)
				h.AssertNil(t, err)
				file, err := os.ReadFile(filepath.Join(tmpDir,"manifests", ref.Name(), ref.Identifier()))
				h.AssertNil(t, err)
				err = json.Unmarshal(file, &manifestList)
				h.AssertNil(t, err)
				for _, hash := range hashes {
					img, err := manifestList.Index.Image(hash)
					h.AssertNotNil(t, err)
					h.AssertNil(t, img)
				}
			})
		})
		when("when --os-version flags passed", func() {
			it("should not return an error when os flag is not specified", func() {
				var node = "traefik:v3.0.0-beta4-windowsservercore-1809"
				command.SetArgs([]string{manifestListName, node, "--os-version", "10.0.14393.1066"})
				err := command.Execute()
				h.AssertNil(t, err)
			})
			it("should include the manifest of given os-version from the image index", func() {
				var manifestList fakeBadIndexStruct
				var osVersion = "10.0.14393.1066"
				var node = "traefik:v3.0.0-beta4-windowsservercore-1809"
				command.SetArgs([]string{manifestListName, node, "--os-version", osVersion})
				err := command.Execute()
				h.AssertNil(t, err)
				ref, err := name.ParseReference(manifestListName)
				h.AssertNil(t, err)
				file, err := os.ReadFile(filepath.Join(tmpDir,"manifests", ref.Name(), ref.Identifier()))
				h.AssertNil(t, err)
				json.Unmarshal(file, &manifestList)
				indexManifest, err := manifestList.Index.IndexManifest()
				h.AssertNil(t, err)
				h.AssertEq(t, indexManifest.Subject.Platform.OSVersion, osVersion)
			})
		})
		// --feature is reserved for future use
		when("when --features flags passed", func() {
			it("should set expected features", func() {
				var manifestList fakeBadIndexStruct
				var features = []string{
					"someFeature",
				}
				var node = "traefik:v3.0.0-beta4-windowsservercore-1809"
				command.SetArgs([]string{manifestListName, node, "--features", features[0]})
				err := command.Execute()
				h.AssertNil(t, err)
				ref, err := name.ParseReference(manifestListName)
				h.AssertNil(t, err)
				file, err := os.ReadFile(filepath.Join(tmpDir,"manifests", ref.Name(), ref.Identifier()))
				h.AssertNil(t, err)
				json.Unmarshal(file, &manifestList)
				indexManifest, err := manifestList.Index.IndexManifest()
				h.AssertNil(t, err)
				h.AssertContains(t, indexManifest.Subject.Platform.Features, features[0])
			})
		})
		when("when --os-features flags passed", func() {
			it("should set expected os-featres", func() {
				var manifestList fakeBadIndexStruct
				var osFeatures = []string{
					"win32k.sys",
				}
				var node = "traefik:v3.0.0-beta4-windowsservercore-1809"
				command.SetArgs([]string{manifestListName, node, "--os-features", osFeatures[0]})
				err := command.Execute()
				h.AssertNil(t, err)
				ref, err := name.ParseReference(manifestListName)
				h.AssertNil(t, err)
				file, err := os.ReadFile(filepath.Join(tmpDir,"manifests", ref.Name(), ref.Identifier()))
				h.AssertNil(t, err)
				json.Unmarshal(file, &manifestList)
				indexManifest, err := manifestList.Index.IndexManifest()
				h.AssertNil(t, err)
				h.AssertContains(t, indexManifest.Subject.Platform.OSFeatures, osFeatures[0])
			})
		})
		when("when --annotations flags passed", func() {
			it("should accept annotations", func() {
				var manifestList fakeBadIndexStruct
				var annotations = []string{
					"com.example.key1=value1",
					"com.example.key2=value2",
				}
				var node = "traefik:v3.0.0-beta4-windowsservercore-1809"
				command.SetArgs([]string{manifestListName, node, "--annotations", strings.Join(annotations, ";")})
				err := command.Execute()
				h.AssertNil(t, err)
				ref, err := name.ParseReference(manifestListName)
				h.AssertNil(t, err)
				file, err := os.ReadFile(filepath.Join(tmpDir,"manifests", ref.Name(), ref.Identifier()))
				h.AssertNil(t, err)
				json.Unmarshal(file, &manifestList)
				indexManifest, err := manifestList.Index.IndexManifest()
				h.AssertNil(t, err)
				h.AssertEq(t, h.MapToStringSlice(indexManifest.Subject.Annotations), annotations)
			})
		})
		when("should throw an error when", func() {
			it("has no args passed", func() {
				command.SetArgs([]string{})
				err := command.Execute()
				h.AssertNotNil(t, err)
			})
			it("has manifest list reference is incorrect", func() {
				var node = "traefikWr0!/\ng`"
				command.SetArgs([]string{node})
				err := command.Execute()
				h.AssertNotNil(t, err)
			})
			it("has  manifest reference is incorrect", func() {
				var node = "traefik:v3.0.0-beta4-windowsservercore-1809"
				command.SetArgs([]string{manifestListName, node+"!#@:'"})
				err := command.Execute()
				h.AssertNotNil(t, err)
			})
			it("has manifest passed in-place of manifest list on first arg", func() {
				var node = "alpine@sha256:777e2106170c66742ddbe77f703badb7dc94d9a5b1dc2c4a01538fad9aef56bb"
				command.SetArgs([]string{node, manifestListName})
				err := command.Execute()
				h.AssertNotNil(t, err)
			})
		})
		when("should warn when",func() {
			it("manifest is passed on second arg with --all option", func() {
				var node = "alpine@sha256:777e2106170c66742ddbe77f703badb7dc94d9a5b1dc2c4a01538fad9aef56bb"
				command.SetArgs([]string{manifestListName, node, "--all"})
				err := command.Execute()
				h.AssertNil(t, err)
				h.AssertNotNil(t, outBuf)
			})
		})
		when("manifest/ImageIndex", func() {
			it("manifest list locally available", func() {
				var node = "alpine@sha256:777e2106170c66742ddbe77f703badb7dc94d9a5b1dc2c4a01538fad9aef56bb"
				command.SetArgs([]string{manifestListName, node, "--all"})
				ref, err := name.ParseReference(manifestListName)
				h.AssertNil(t, err)
				var path = filepath.Join(tmpDir,"manifests", ref.Name(), ref.Identifier())
				_, err = os.Stat(path)
				h.AssertNil(t, err)
				err = command.Execute()
				h.AssertNil(t, err)
				h.AssertNotNil(t, outBuf)
				_, err = os.Stat(path)
				h.AssertNil(t, err)
			})
			it("manifest list available at registry only", func() {
				var node = "alpine@sha256:777e2106170c66742ddbe77f703badb7dc94d9a5b1dc2c4a01538fad9aef56bb"
				ref, err := name.ParseReference(node)
				h.AssertNil(t, err)
				var path = filepath.Join(tmpDir,"manifests", ref.Name(), ref.Identifier())
				_, err = os.Stat(path)
				h.AssertNotNil(t, err)
				command.SetArgs([]string{"alpine:3.18", node, "--all"})
				err = command.Execute()
				h.AssertNil(t, err)
				_, err = os.Stat(path)
				h.AssertNil(t, err)
			})
		})
	})
}
