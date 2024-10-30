package add

import (
	"fmt"
	"strings"
)

var _ AddCMDStringer = (*addCMDStringer)(nil)

func (as addCMDStringer) String() string {
	return as.
		KeepGitDir().
		Checksum().
		Chown().
		Chmod().
		Link().
		Parent().
		Exclude().
		Sources().
		Dest().
		toString()
}

func (as *addCMDStringer) Sources() *addCMDStringer {
	s := strings.Join(as.sources, " ")
	if strings.HasSuffix(as.builder.String(), " ") {
		as.builder.WriteString(s)
	} else {
		as.builder.WriteString(" " + s)
	}

	return as
}


func (as *addCMDStringer) Dest() *addCMDStringer {
	if strings.HasSuffix(as.builder.String(), " ") {
		as.builder.WriteString(as.dest)
	} else {
		as.builder.WriteString(" " + as.dest)
	}

	return as
}

func (as *addCMDStringer) Link() *addCMDStringer {
	if strings.HasSuffix(as.builder.String(), " ") {
		as.builder.WriteString(
			as.options.Link.FormatFlag(),
		)
	} else {
		as.builder.WriteString(
			" " + as.options.Link.FormatFlag(),
		)
	}

	return as
}

func (as *addCMDStringer) Parent() *addCMDStringer {
	if strings.HasSuffix(as.builder.String(), " ") {
		as.builder.WriteString(
			as.options.Parent.FormatFlag(),
		)
	} else {
		as.builder.WriteString(
			" " + as.options.Parent.FormatFlag(),
		)
	}

	return as
}

func (as *addCMDStringer) Exclude() *addCMDStringer {
	if strings.HasSuffix(as.builder.String(), " ") {
		as.builder.WriteString(
			as.options.Exclude.FormatFlag(),
		)
	} else {
		as.builder.WriteString(
			" " + as.options.Exclude.FormatFlag(),
		)
	}

	return as
}

func (as *addCMDStringer) Chown() *addCMDStringer {
	if strings.HasSuffix(as.builder.String(), " ") {
		as.builder.WriteString(as.options.Chown.FormatFlag())
	} else {
		as.builder.WriteString(" " + as.options.Chown.FormatFlag())
	}

	return as
}

func (as *addCMDStringer) Chmod() *addCMDStringer {
	if as.options.Chmod != 0 {
		if strings.HasSuffix(as.builder.String(), " ") {
			as.builder.WriteString(as.options.Chmod.FormatFlag())
		} else {
			as.builder.WriteString(" " + as.options.Chmod.FormatFlag())
		}
	}

	return as
}

func (as *addCMDStringer) Checksum() *addCMDStringer {
	if strings.HasSuffix(as.builder.String(), " ") {
		as.builder.WriteString(as.options.Checksum.FormatFlag())
	} else {
		as.builder.WriteString(" " + as.options.Checksum.FormatFlag())
	}

	return as
}

func (as *addCMDStringer) KeepGitDir() *addCMDStringer {
	if strings.HasSuffix(as.builder.String(), " ") {
		as.builder.WriteString(as.options.KeepGitDir.FormatFlag())
	} else {
		as.builder.WriteString(" " + as.options.KeepGitDir.FormatFlag())
	}

	return as
}

func (as *addCMDStringer) toString() string {
	return fmt.Sprintf("%s # buildpack", as.builder.String())
}
