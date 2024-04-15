package dist

import (
	"errors"
	"runtime"

	v1 "github.com/google/go-containerregistry/pkg/v1"
)

type Target struct {
	OS            string         `json:"os" toml:"os"`
	Arch          string         `json:"arch" toml:"arch"`
	ArchVariant   string         `json:"variant,omitempty" toml:"variant,omitempty"`
	Distributions []Distribution `json:"distributions,omitempty" toml:"distributions,omitempty"`
	Specs         TargetSpecs    `json:"specs,omitempty" toml:"specs,omitempty"`
}

type Distribution struct {
	Name     string   `json:"name,omitempty" toml:"name,omitempty"`
	Versions []string `json:"versions,omitempty" toml:"versions,omitempty"`
}

type TargetSpecs struct {
	Features       []string          `json:"features,omitempty" toml:"features,omitempty"`
	OSFeatures     []string          `json:"os.features,omitempty" toml:"os.features,omitempty"`
	URLs           []string          `json:"urls,omitempty" toml:"urls,omitempty"`
	Annotations    map[string]string `json:"annotations,omitempty" toml:"annotations,omitempty"`
	Flatten        bool              `json:"flatten,omitempty" toml:"flatten,omitempty"`
	FlattenExclude []string          `json:"flatten.exclude,omitempty" toml:"flatten.exclude,omitempty"`
	Labels         map[string]string `json:"labels,omitempty" toml:"labels,omitempty"`
	OSVersion      string            `json:"os.version,omitempty" toml:"os.version,omitempty"`
	Path           string            `json:"path,omitempty" toml:"path,omitempty"`
}

func (t Target) Range(op func(target Target, distroName, distroVersion string) error) error {
	emptyString := ""
	if len(t.Distributions) == 0 {
		return op(t, emptyString, emptyString)
	}
	for _, distro := range t.Distributions {
		if len(distro.Versions) == 0 {
			if err := op(t, distro.Name, emptyString); err != nil {
				return err
			}
			continue
		}
		for _, version := range distro.Versions {
			if err := op(t, distro.Name, version); err != nil {
				return err
			}
		}
	}
	return nil
}

func (t *Target) MultiArch() bool {
	length := 0
	t.Range(func(target Target, distroName, distroVersion string) error {
		length++
		return nil
	})

	return length > 1
}

func (t *Target) Platform() *v1.Platform {
	return &v1.Platform{
		OS:           t.OS,
		Architecture: t.Arch,
		Variant:      t.ArchVariant,
		OSVersion:    t.Specs.OSVersion,
		Features:     t.Specs.Features,
		OSFeatures:   t.Specs.OSFeatures,
	}
}

func (t *Target) Annotations() (map[string]string, error) {
	if len(t.Distributions) == 0 {
		return nil, errors.New("unable to get annotations: distroless target provided")
	}

	distro := t.Distributions[0]
	return map[string]string{
		"io.buildpacks.base.distro.name":    distro.Name,
		"io.buildpacks.base.distro.version": distro.Versions[0],
	}, nil
}

func (t *Target) URLs() []string {
	return t.Specs.URLs
}

func DefaultTarget() Target {
	return Target{
		OS:   runtime.GOOS,
		Arch: runtime.GOARCH,
	}
}
