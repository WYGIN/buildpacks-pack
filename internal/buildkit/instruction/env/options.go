package env

import (
	"github.com/buildpacks/pack/internal/buildkit/instruction"
	"github.com/moby/buildkit/client/llb"
)

func WithKeyValue(key, value string) func(*instruction.EnvOp) error {
	return func(ao *instruction.EnvOp) error {
		ao.EnvList = llb.EnvList{}
		ao.EnvList.AddOrReplace(key, value)
		return nil
	}
}

func WithKeyValuePair(kv ...instruction.KeyValuePair) func(*instruction.EnvOp) error {
	return func(ao *instruction.EnvOp) error {
		ao.EnvList = llb.EnvList{}
		for _, kvPair := range kv {
			ao.EnvList.AddOrReplace(kvPair.Key, kvPair.Value)
		}

		return nil
	}
}
