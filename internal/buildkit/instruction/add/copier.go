package add

import (
	"net/url"
	"path"
	"strings"

	"github.com/buildpacks/pack/internal/buildkit/instruction"
	"github.com/moby/buildkit/client/llb"
	"github.com/moby/buildkit/util/gitutil"
	"github.com/moby/buildkit/util/system"
	"github.com/pkg/errors"
)

func NewSourceCopier(a add, src SourceAndDest) SourceCopier {
	var op = SourceCopierOps{
		src: src,
		a: a,
	}

	if instruction.IsGitSource(src.src) {
		return gitSourceCopier{SourceCopierOps: op}
	}

	if instruction.IsHttpSource(src.src) {
		return httpSourceCopier{SourceCopierOps: op}
	}

	return localSourceCopier{SourceCopierOps: op}
}

// Copy implements Copier.
func (g gitSourceCopier) Copy(fileAction FileAction) error {
	var (
		src = g.src.src
		dest = g.src.dest
		add = g.a
	)

	// no need to concern about error
	// it is already handled in [instruction.IsGitSource]
	ref, _ := gitutil.ParseGitRef(src)
	commit := ref.Commit
	if ref.SubDir != "" {
		commit += ":" + ref.SubDir
	}

	var gitOps = []llb.GitOption{/*llb.WithCustomName(pgName)*/}
	if g.a.options.KeepGitDir {
		gitOps = append(gitOps, llb.KeepGitDir())
	}

	st := llb.Git(ref.Remote, commit, gitOps...)
	opts := append([]llb.CopyOption{&llb.CopyInfo{
		CreateDestPath: true,
	}}, add.copyOps(src))

	fileAction(st, "/", dest, opts...)
	return nil
}

// Copy implements Copier.
func (g httpSourceCopier) Copy(fileAction FileAction) error {
	var (
		src = g.src.src
		dest = g.src.dest
		add = g.a
	)

	u, err := url.Parse(src)
	f := "__unnamed__"

	if err == nil {
		if base := path.Base(u.Path); base != "." && base != "/" {
			f = base
		}
	}

	st := llb.HTTP(src, llb.Filename(f), /*llb.WithCustomName(pgName),*/ llb.Checksum(add.options.Digest), /*dfCmd(cfg.params)*/)
	opts := append([]llb.CopyOption{&llb.CopyInfo{
		CreateDestPath: true,
	}}, add.copyOps(src))

	fileAction(st, f, dest, opts...)
	return nil
}

// Copy implements Copier.
func (g localSourceCopier) Copy(fileAction FileAction) (err error) {
	var (
		src = g.src.src
		dest = g.src.dest
		add = g.a
	)

	var patterns []string
	if add.options.Parent {
		src, patterns, err = normalizeParentPath(add.platform.OS, src, patterns)
	}

	if err != nil {
		return err
	}

	src, err = system.NormalizePath("/", src, add.platform.OS, false)
	if err != nil {
		return errors.Wrap(err, "removing drive letter")
	}

	ops := append([]llb.CopyOption{&llb.CopyInfo{
		FollowSymlinks:      true,
		CopyDirContentsOnly: true,
		IncludePatterns:     patterns,
		AttemptUnpack:       isAddCommand(add),
		CreateDestPath:      true,
		AllowWildcard:       true,
		AllowEmptyWildcard:  true,
	}}, add.copyOps(src))

	fileAction(add.state, src, dest, ops...)
	return nil
}

func normalizeParentPath(os, src string, patterns []string) (string, []string, error) {
	// detect optional pivot point
	parent, pattern, ok := strings.Cut(src, "/./")
	if !ok {
		pattern = src
		src = "/"
	} else {
		src = parent
	}

	pattern, err := system.NormalizePath("/", pattern, os, false)
	if err != nil {
		return src, patterns, errors.Wrap(err, "removing drive letter")
	}

	patterns = []string{strings.TrimPrefix(pattern, "/")}

	return src, patterns, nil
}

func isAddCommand(a add) bool {
	switch {
	case a.options.Checksum.String() != "",
		instruction.IsGitSource(a.sources[0]),
		instruction.IsHttpSource(a.sources[0]),
		a.options.attemptUnpack:
			return true
	}

	return false
}