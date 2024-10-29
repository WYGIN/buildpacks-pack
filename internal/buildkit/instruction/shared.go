package instruction

import (
	"strings"

	"github.com/moby/buildkit/util/gitutil"
)

func IsHttpSource(src string) bool {
	if !strings.HasPrefix(src, "http://") && !strings.HasPrefix(src, "https://") {
		return false
	}
	// https://github.com/ORG/REPO.git is a git source, not an http source
	if gitRef, gitErr := gitutil.ParseGitRef(src); gitRef != nil && gitErr == nil {
		return false
	}
	return true
}

func IsGitSource(src string) bool {
	ref, err := gitutil.ParseGitRef(src)

	return err == nil && ref.IndistinguishableFromLocal
}