package diff

import (
	"strings"

	"github.com/pkg/errors"
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
	lines := strings.Split(string(b), "\n")

	var cmds []*cmd
	for _, c := range d.Changes {
		cmds = append(cmds, loadCmd(c, false))
	}
	newLines, err := applyCmds(lines, cmds...)
	if err != nil {
		return nil, errors.Wrap(err, "unable to apply commands")
	}

	return []byte(strings.Join(newLines, "\n")), nil
}

func ApplyInverse(b []byte, d *Diff) ([]byte, error) {
	lines := strings.Split(string(b), "\n")

	var cmds []*cmd
	for _, c := range d.Changes {
		cmds = append(cmds, loadCmd(c, true))
	}
	newLines, err := applyCmds(lines, cmds...)
	if err != nil {
		return nil, errors.Wrap(err, "unable to apply commands")
	}

	return []byte(strings.Join(newLines, "\n")), nil
}

func applyCmds(lines []string, cmds ...*cmd) ([]string, error) {
	ret := make([]string, 0, len(lines))
	currIdx := 0

	for _, c := range cmds {
		for ; currIdx <= c.From + 1; currIdx++ {
			ret = append(ret, lines[currIdx])
		}
		for _, a := range c.Added {
			ret = append(ret, a)
		}
		currIdx += len(c.Deleted)
	}
	for ; currIdx < len(lines); currIdx++ {
		ret = append(ret, lines[currIdx])
	}
	return ret, nil
}

func loadCmd(c *Change, inverse bool) *cmd {
	x := &cmd{From: c.From, To: c.To}
	for _, l := range c.Lines {
		lv := strings.TrimSuffix(l.V, "\n")
		switch l.T {
		case "context":
			if len(x.Deleted) == 0 && len(x.Added) == 0 {
				x.PreCtx = append(x.PreCtx, lv)
			} else {
				x.PostCtx = append(x.PostCtx, lv)
			}
		case "deleted":
			x.Deleted = append(x.Deleted, lv)
		case "added":
			x.Added = append(x.Added, lv)
		}
	}
	if inverse {
		a := x.Added
		x.Added = x.Deleted
		x.Deleted = a
	}
	return x
}
