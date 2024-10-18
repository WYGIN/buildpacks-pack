package label

import "github.com/buildpacks/pack/internal/buildkit/instruction"

func WithLabels(labels ...instruction.KeyValuePair) func(*instruction.LabelOp) error {
	return func(lo *instruction.LabelOp) error {
		lo.Labels = labels
		return nil
	}
}