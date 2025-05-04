package diff

import (
	"strings"

	"github.com/pkg/errors"
	"github.com/samber/lo"

	"projectforge.dev/projectforge/app/util"
)

type cmd struct {
	From    int
	To      int
	PreCtx  []string
	Deleted []string
	Added   []string
	PostCtx []string
}

func Apply(b []byte, d *Diff) ([]byte, error) {
	lines := util.StringSplitLines(string(b))
	cmds := lo.Map(d.Changes, func(c *Change, _ int) *cmd {
		return loadCmd(c, false)
	})
	newLines, err := applyCmds(lines, cmds...)
	if err != nil {
		return nil, errors.Wrap(err, "unable to apply commands")
	}

	return []byte(util.StringJoin(newLines, util.StringDetectLinebreak(string(b)))), nil
}

func ApplyInverse(b []byte, d *Diff) ([]byte, error) {
	lines := util.StringSplitLines(string(b))
	cmds := lo.Map(d.Changes, func(c *Change, _ int) *cmd {
		return loadCmd(c, true)
	})
	newLines, err := applyCmds(lines, cmds...)
	if err != nil {
		return nil, errors.Wrap(err, "unable to apply commands")
	}

	return []byte(util.StringJoin(newLines, util.StringDetectLinebreak(string(b)))), nil
}

func applyCmds(lines []string, cmds ...*cmd) ([]string, error) {
	ret := util.NewStringSlice(make([]string, 0, len(lines)))
	var currIdx int
	lo.ForEach(cmds, func(c *cmd, _ int) {
		for ; currIdx <= c.From+1; currIdx++ {
			ret.Push(lines[currIdx])
		}
		ret.Push(c.Added...)
		currIdx += len(c.Deleted)
	})
	for ; currIdx < len(lines); currIdx++ {
		ret.Push(lines[currIdx])
	}
	return ret.Slice, nil
}

func loadCmd(c *Change, inverse bool) *cmd {
	x := &cmd{From: c.From, To: c.To}
	lo.ForEach(c.Lines, func(l *Line, _ int) {
		lv := strings.TrimSuffix(strings.TrimSuffix(l.V, "\n"), "\r")
		switch l.T {
		case contextKey:
			if len(x.Deleted) == 0 && len(x.Added) == 0 {
				x.PreCtx = append(x.PreCtx, lv)
			} else {
				x.PostCtx = append(x.PostCtx, lv)
			}
		case deletedKey:
			x.Deleted = append(x.Deleted, lv)
		case addedKey:
			x.Added = append(x.Added, lv)
		}
	})
	if inverse {
		x.Added, x.Deleted = x.Deleted, x.Added
	}
	return x
}
