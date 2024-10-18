package arg

import (
	"github.com/buildpacks/pack/internal/buildkit/instruction"
	"github.com/moby/buildkit/client/llb"
)

func WithKeyValue(key, value string) func(*instruction.ArgOp) error {
	return func(ao *instruction.ArgOp) error {
		ao.Args = make(llb.EnvList, 1)
		ao.Args.AddOrReplace(key, value)
		return nil
	}
}

func WithKeyValuePair(kv ...instruction.KeyValuePair) func(*instruction.ArgOp) error {
	return func(ao *instruction.ArgOp) error {
		ao.Args = make(llb.EnvList, 0, len(kv))
		for _, kvPair := range kv {
			ao.Args.AddOrReplace(kvPair.Key, kvPair.Value)
		}

		return nil
	}
}
