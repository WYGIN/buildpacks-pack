package stage

import "github.com/moby/buildkit/client/llb"

type Stage interface {
	Build() (llb.State, error)
}

type Remote interface {}

type Git interface {}

type Local interface {}

type Http interface {}