package add

import "github.com/buildpacks/pack/internal/buildkit/instruction"

var _ AddHistroryCommiter = (*addHistroryCommiter)(nil)

func NewAddCommitHistoryBuilder() instruction.HistoryCommiter {
	return nil
}

func (ahc addHistroryCommiter) Commit() error {
	_ = ahc.add
	_ = ahc.frontend
	return nil
}