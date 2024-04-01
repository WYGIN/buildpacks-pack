package builder

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/BurntSushi/toml"
	"github.com/pkg/errors"

	"github.com/buildpacks/pack/buildpackage"
	"github.com/buildpacks/pack/internal/config"
	"github.com/buildpacks/pack/internal/style"
	"github.com/buildpacks/pack/pkg/dist"
)

// Config is a builder configuration file
type Config struct {
	Description     string           `toml:"description"`
	Buildpacks      ModuleCollection `toml:"buildpacks"`
	Extensions      ModuleCollection `toml:"extensions"`
	Order           dist.Order       `toml:"order"`
	OrderExtensions dist.Order       `toml:"order-extensions"`
	Stack           StackConfig      `toml:"stack"`
	Lifecycle       LifecycleConfig  `toml:"lifecycle"`
	Run             RunConfig        `toml:"run"`
	Build           BuildConfig      `toml:"build"`
	WithTargets     []dist.Target    `toml:"targets,omitempty"`
}

type MultiArchConfig struct {
	Config
	flagTargets     []dist.Target
	relativeBaseDir string
}

// ModuleCollection is a list of ModuleConfigs
type ModuleCollection []ModuleConfig

// ModuleConfig details the configuration of a Buildpack or Extension
type ModuleConfig struct {
	dist.ModuleInfo
	dist.ImageOrURI
}

func (c *ModuleConfig) DisplayString() string {
	if c.ModuleInfo.FullName() != "" {
		return c.ModuleInfo.FullName()
	}

	return c.ImageOrURI.DisplayString()
}

// StackConfig details the configuration of a Stack
type StackConfig struct {
	ID              string   `toml:"id"`
	BuildImage      string   `toml:"build-image"`
	RunImage        string   `toml:"run-image"`
	RunImageMirrors []string `toml:"run-image-mirrors,omitempty"`
}

// LifecycleConfig details the configuration of the Lifecycle
type LifecycleConfig struct {
	URI     string `toml:"uri"`
	Version string `toml:"version"`
}

// RunConfig set of run image configuration
type RunConfig struct {
	Images []RunImageConfig `toml:"images"`
}

// RunImageConfig run image id and mirrors
type RunImageConfig struct {
	Image   string   `toml:"image"`
	Mirrors []string `toml:"mirrors,omitempty"`
}

// BuildConfig build image configuration
type BuildConfig struct {
	Image string           `toml:"image"`
	Env   []BuildConfigEnv `toml:"env"`
}

type Suffix string

const (
	NONE     Suffix = ""
	DEFAULT  Suffix = "default"
	OVERRIDE Suffix = "override"
	APPEND   Suffix = "append"
	PREPEND  Suffix = "prepend"
)

type BuildConfigEnv struct {
	Name   string `toml:"name"`
	Value  string `toml:"value"`
	Suffix Suffix `toml:"suffix,omitempty"`
	Delim  string `toml:"delim,omitempty"`
}

// ReadConfig reads a builder configuration from the file path provided and returns the
// configuration along with any warnings encountered while parsing
func ReadConfig(path string) (config Config, warnings []string, err error) {
	file, err := os.Open(filepath.Clean(path))
	if err != nil {
		return Config{}, nil, errors.Wrap(err, "opening config file")
	}
	defer file.Close()

	config, err = parseConfig(file)
	if err != nil {
		return Config{}, nil, errors.Wrapf(err, "parse contents of '%s'", path)
	}

	if len(config.Order) == 0 {
		warnings = append(warnings, fmt.Sprintf("empty %s definition", style.Symbol("order")))
	}

	config.mergeStackWithImages()

	return config, warnings, nil
}

func (c *MultiArchConfig) Targets() []dist.Target {
	if len(c.flagTargets) != 0 {
		return c.flagTargets
	}

	return c.WithTargets
}

func ReadMultiArchConfig(path string, flagTargets []dist.Target) (config MultiArchConfig, warnings []string, err error) {
	file, err := os.Open(filepath.Clean(path))
	if err != nil {
		return MultiArchConfig{}, nil, errors.Wrap(err, "opening config file")
	}
	defer file.Close()

	config, err = parseMultiArchConfig(file)
	if err != nil {
		return MultiArchConfig{}, nil, errors.Wrapf(err, "parse contents of '%s'", path)
	}

	if len(config.Order) == 0 {
		warnings = append(warnings, fmt.Sprintf("empty %s definition", style.Symbol("order")))
	}

	config.mergeStackWithImages()
	if len(flagTargets) != 0 {
		config.flagTargets = flagTargets
	}
	return config, warnings, nil
}

func (c *MultiArchConfig) BuilderConfigs(getIndexManifest buildpackage.GetIndexManifestFn) (configs []Config, err error) {
	for _, target := range c.Targets() {
		if err := target.Range(func(target dist.Target, distroName, distroVersion string) error {
			cfg, err := c.processTarget(target, "", "", getIndexManifest)
			configs = append(configs, cfg)
			return err
		}); err != nil {
			return configs, err
		}
	}
	return configs, nil
}

func copyConfig(config Config) Config {
	return Config{
		Description:     config.Description,
		Buildpacks:      make(ModuleCollection, len(config.Buildpacks)),
		Extensions:      make(ModuleCollection, len(config.Extensions)),
		Order:           make(dist.Order, len(config.Order)),
		OrderExtensions: make(dist.Order, len(config.OrderExtensions)),
		Stack: StackConfig{
			RunImageMirrors: make([]string, len(config.Stack.RunImageMirrors)),
		},
		Run: RunConfig{
			Images: make([]RunImageConfig, len(config.Run.Images)),
		},
	}
}

func (c *MultiArchConfig) processTarget(target dist.Target, distroName, distroVersion string, getIndexManifest buildpackage.GetIndexManifestFn) (config Config, err error) {
	config = copyConfig(c.Config)
	target = buildpackage.ProcessTarget(target, distroName, distroVersion)
	for i, bp := range c.Config.Buildpacks {
		if bp.URI != "" {
			if config.Buildpacks[i].URI, err = buildpackage.GetRelativeURI(bp.URI, c.relativeBaseDir, &target, getIndexManifest); err != nil {
				return config, err
			}
		}
	}

	for i, ext := range c.Config.Extensions {
		if ext.URI != "" {
			if config.Extensions[i].URI, err = buildpackage.GetRelativeURI(ext.URI, c.relativeBaseDir, &target, getIndexManifest); err != nil {
				return config, err
			}
		}
	}

	if img := c.Config.Build.Image; img != "" {
		if config.Build.Image, err = buildpackage.ParseURItoString(img, target, getIndexManifest); err != nil {
			return config, err
		}
	}

	for i, runImg := range c.Config.Run.Images {
		config.Run.Images[i].Image, err = buildpackage.ParseURItoString(runImg.Image, target, getIndexManifest)
		if len(config.Run.Images[i].Mirrors) == 0 {
			config.Run.Images[i].Mirrors = make([]string, len(runImg.Mirrors))
		}

		if err != nil {
			for j, mirror := range runImg.Mirrors {
				if config.Run.Images[i].Mirrors[j], err = buildpackage.ParseURItoString(mirror, target, getIndexManifest); err == nil {
					break
				}
			}

			if err != nil {
				return config, err
			}
		}
	}

	if img := c.Config.Stack.BuildImage; img != "" {
		if config.Stack.BuildImage, err = buildpackage.ParseURItoString(img, target, getIndexManifest); err != nil {
			return config, err
		}
	}

	if img := c.Config.Stack.RunImage; img != "" {
		config.Stack.RunImage, err = buildpackage.ParseURItoString(img, target, getIndexManifest)
	}

	if err != nil {
		for i, mirror := range c.Config.Stack.RunImageMirrors {
			if config.Stack.RunImageMirrors[i], err = buildpackage.ParseURItoString(mirror, target, getIndexManifest); err == nil {
				break
			}
		}
	}

	config.Order = c.Order
	config.OrderExtensions = c.OrderExtensions
	config.WithTargets = []dist.Target{target}
	return config, err
}

func (c *MultiArchConfig) MultiArch() bool {
	targets := c.Targets()
	if len(targets) > 1 {
		return true
	}

	targetsLen := 0
	for _, target := range targets {
		target.Range(func(_ dist.Target, _, _ string) error {
			targetsLen++
			return nil
		})
	}

	return targetsLen > 1
}

// ValidateConfig validates the config
func ValidateConfig(c Config) error {
	if c.Build.Image == "" && c.Stack.BuildImage == "" {
		return errors.New("build.image is required")
	} else if c.Build.Image != "" && c.Stack.BuildImage != "" && c.Build.Image != c.Stack.BuildImage {
		return errors.New("build.image and stack.build-image do not match")
	}

	if len(c.Run.Images) == 0 && (c.Stack.RunImage == "" || c.Stack.ID == "") {
		return errors.New("run.images are required")
	}

	for _, runImage := range c.Run.Images {
		if runImage.Image == "" {
			return errors.New("run.images.image is required")
		}
	}

	if c.Stack.RunImage != "" && c.Run.Images[0].Image != c.Stack.RunImage {
		return errors.New("run.images and stack.run-image do not match")
	}

	return nil
}

func (c *Config) mergeStackWithImages() {
	// RFC-0096
	if c.Build.Image != "" {
		c.Stack.BuildImage = c.Build.Image
	} else if c.Build.Image == "" && c.Stack.BuildImage != "" {
		c.Build.Image = c.Stack.BuildImage
	}

	if len(c.Run.Images) != 0 {
		// use the first run image as the "stack"
		c.Stack.RunImage = c.Run.Images[0].Image
		c.Stack.RunImageMirrors = c.Run.Images[0].Mirrors
	} else if len(c.Run.Images) == 0 && c.Stack.RunImage != "" {
		c.Run.Images = []RunImageConfig{{
			Image:   c.Stack.RunImage,
			Mirrors: c.Stack.RunImageMirrors,
		},
		}
	}
}

// parseConfig reads a builder configuration from file
func parseConfig(file *os.File) (Config, error) {
	builderConfig := Config{}
	tomlMetadata, err := toml.NewDecoder(file).Decode(&builderConfig)
	if err != nil {
		return Config{}, errors.Wrap(err, "decoding toml contents")
	}

	undecodedKeys := tomlMetadata.Undecoded()
	if len(undecodedKeys) > 0 {
		unknownElementsMsg := config.FormatUndecodedKeys(undecodedKeys)

		return Config{}, errors.Errorf("%s in %s",
			unknownElementsMsg,
			style.Symbol(file.Name()),
		)
	}

	return builderConfig, nil
}

// parseMultiArchConfig reads a builder configuration from file
func parseMultiArchConfig(file *os.File) (MultiArchConfig, error) {
	multiArchBuilderConfig := MultiArchConfig{}
	tomlMetadata, err := toml.NewDecoder(file).Decode(&multiArchBuilderConfig)
	if err != nil {
		return MultiArchConfig{}, errors.Wrap(err, "decoding MultiArchBuilder Toml")
	}

	undecodedKeys := tomlMetadata.Undecoded()
	if len(undecodedKeys) > 0 {
		unknownElementsMsg := config.FormatUndecodedKeys(undecodedKeys)

		return MultiArchConfig{}, errors.Errorf("%s in %s",
			unknownElementsMsg,
			style.Symbol(file.Name()),
		)
	}

	return multiArchBuilderConfig, nil
}

func ParseBuildConfigEnv(env []BuildConfigEnv, path string) (envMap map[string]string, warnings []string, err error) {
	envMap = map[string]string{}
	var appendOrPrependWithoutDelim = 0
	for _, v := range env {
		if name := v.Name; name == "" {
			return nil, nil, errors.Wrapf(errors.Errorf("env name should not be empty"), "parse contents of '%s'", path)
		}
		if val := v.Value; val == "" {
			warnings = append(warnings, fmt.Sprintf("empty value for key/name %s", style.Symbol(v.Name)))
		}
		suffixName, delimName, err := getBuildConfigEnvFileName(v)
		if err != nil {
			return envMap, warnings, err
		}
		if val, ok := envMap[suffixName]; ok {
			warnings = append(warnings, fmt.Sprintf(errors.Errorf("overriding env with name: %s and suffix: %s from %s to %s", style.Symbol(v.Name), style.Symbol(string(v.Suffix)), style.Symbol(val), style.Symbol(v.Value)).Error(), "parse contents of '%s'", path))
		}
		if val, ok := envMap[delimName]; ok {
			warnings = append(warnings, fmt.Sprintf(errors.Errorf("overriding env with name: %s and delim: %s from %s to %s", style.Symbol(v.Name), style.Symbol(v.Delim), style.Symbol(val), style.Symbol(v.Value)).Error(), "parse contents of '%s'", path))
		}
		if delim := v.Delim; delim != "" && delimName != "" {
			envMap[delimName] = delim
		}
		envMap[suffixName] = v.Value
	}

	for k := range envMap {
		name, suffix, err := getFilePrefixSuffix(k)
		if err != nil {
			continue
		}
		if _, ok := envMap[name+".delim"]; (suffix == "append" || suffix == "prepend") && !ok {
			warnings = append(warnings, fmt.Sprintf(errors.Errorf("env with name/key %s with suffix %s must to have a %s value", style.Symbol(name), style.Symbol(suffix), style.Symbol("delim")).Error(), "parse contents of '%s'", path))
			appendOrPrependWithoutDelim++
		}
	}
	if appendOrPrependWithoutDelim > 0 {
		return envMap, warnings, errors.Errorf("error parsing [[build.env]] in file '%s'", path)
	}
	return envMap, warnings, err
}

func getBuildConfigEnvFileName(env BuildConfigEnv) (suffixName, delimName string, err error) {
	suffix, err := getActionType(env.Suffix)
	if err != nil {
		return suffixName, delimName, err
	}
	if suffix == "" {
		suffixName = env.Name
	} else {
		suffixName = env.Name + suffix
	}
	if delim := env.Delim; delim != "" {
		delimName = env.Name + ".delim"
	}
	return suffixName, delimName, err
}

func getActionType(suffix Suffix) (suffixString string, err error) {
	const delim = "."
	switch suffix {
	case NONE:
		return "", nil
	case DEFAULT:
		return delim + string(DEFAULT), nil
	case OVERRIDE:
		return delim + string(OVERRIDE), nil
	case APPEND:
		return delim + string(APPEND), nil
	case PREPEND:
		return delim + string(PREPEND), nil
	default:
		return suffixString, errors.Errorf("unknown action type %s", style.Symbol(string(suffix)))
	}
}

func getFilePrefixSuffix(filename string) (prefix, suffix string, err error) {
	val := strings.Split(filename, ".")
	if len(val) <= 1 {
		return val[0], suffix, errors.Errorf("Suffix might be null")
	}
	if len(val) == 2 {
		suffix = val[1]
	} else {
		suffix = strings.Join(val[1:], ".")
	}
	return val[0], suffix, err
}
