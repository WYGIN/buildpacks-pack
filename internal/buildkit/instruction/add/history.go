package add

import "github.com/buildpacks/pack/internal/buildkit/instruction"

var _ AddHistroryCommiter = (*addHistroryCommiter)(nil)

func NewAddHistoryCommitter(add *add) instruction.HistoryCommiter {
	return addHistroryCommiter{
		add: add,
		frontend: "gateway.v0",
	}
}

func (ahc addHistroryCommiter) Commit() error {
	return nil
}