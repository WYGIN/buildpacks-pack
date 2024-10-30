package customname

import v1 "github.com/opencontainers/image-spec/specs-go/v1"

func CurrentCommandIndex(index *uint) CustomNameFactoryOption {
	return func(cn *customName) error {
		cn.currentCmdIndex = index
		return nil
	}
}

func WithPlatform(platform v1.Platform) CustomNameFactoryOption {
	return func(cn *customName) error {
		cn.platform = &platform
		return nil
	}
}

func WithStageName(name string) CustomNameFactoryOption {
	return func(cn *customName) error {
		cn.stageName = name
		return nil
	}
}

func TotalCmdsCount(total *uint) CustomNameFactoryOption {
	return func(cn *customName) error {
		cn.totalCmdsCount = total
		return nil
	}
}

var OnBuildCMD CustomNameFactoryOption = func(cn *customName) error {
	cn.isOnBuildCmd = true
	return nil
}

var IgnoreCmdIndexIncrement CustomNameFactoryOption = func(cn *customName) error {
	cn.ignoreCMDIndexIncrement = true
	return nil
}