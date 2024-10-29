package instruction

import (
	"fmt"
	"strings"
)

func (l Link) FormatFlag() string {
	if l {
		return fmt.Sprintf("--link=%t", l)
	}

	return ""
}

func (p Parent) FormatFlag() string {
	if p {
		return fmt.Sprintf("--parent=%t", p)
	}

	return ""
}

func (e Exclude) FormatFlag() string {
	builder := strings.Builder{}

	var toExcludeFormat = func(pattern string) string {
		return fmt.Sprintf("--exclude=%s", pattern)
	}

	for i, pattern := range e {
		if i == 0 {
			builder.WriteString(toExcludeFormat(pattern))
		} else {
			builder.WriteString(" " + toExcludeFormat(pattern))
		}
	}

	return builder.String()
}

func (c Chown) FormatFlag() string {
	var chown = "--chown="
	var dup = chown

	chown += c.User.Name
	chown +=  fmt.Sprint(c.User.UID)

	if c.Group != nil {
		chown += ":"

		chown += c.Group.Name
		chown += fmt.Sprint(c.Group.UID)
	}

	if dup != chown {
		return chown
	}

	return ""
}

func (c Chmod) FormatFlag() string {
	if c == 0 {
		return ""
	}

	return fmt.Sprintf("--chmod=%o", c)
}

func (c Chmod) IsNil() bool {
	return c == 0o0
}

func (ck Checksum) FormatFlag() string {
	if c := ck.String(); c == "" {
		return c
	}

	return fmt.Sprintf("--checksum=%s", ck.String())
}

func (k KeepGitDir) FormatFlag() string {
	if k {
		return fmt.Sprintf("--keep-git-dir=%t", k)
	}

	return ""
}
