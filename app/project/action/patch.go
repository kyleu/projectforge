package action

import (
	"github.com/samber/lo"

	"projectforge.dev/projectforge/app/util"
)

type PatchSection struct {
	Type  string   `json:"type"`
	Lines []string `json:"lines"`
}

const (
	patchTypeAdded     = "added"
	patchTypeRemoved   = "removed"
	patchTypeUnchanged = "unchanged"
)

func (p *PatchSection) CSSClass() string {
	switch p.Type {
	case patchTypeUnchanged:
		return "color-muted"
	case patchTypeAdded:
		return "success"
	case patchTypeRemoved:
		return "error"
	default:
		return ""
	}
}

func (p *PatchSection) Render() string {
	return util.StringJoin(lo.Map(p.Lines, func(l string, _ int) string {
		if l == "" {
			return ""
		}
		return l[1:]
	}), "\n")
}

func (p *PatchSection) Useful() bool {
	return p.Type == patchTypeAdded || p.Type == patchTypeRemoved
}

type PatchSections []*PatchSection

func (ps PatchSections) AddLine(t string, l string) PatchSections {
	if len(ps) == 0 {
		return append(ps, &PatchSection{Type: t, Lines: []string{l}})
	}
	last := ps[len(ps)-1]
	if last.Type == t {
		last.Lines = append(last.Lines, l)
		return ps
	}
	return append(ps, &PatchSection{Type: t, Lines: []string{l}})
}

func ParsePatchString(s string) PatchSections {
	var ret PatchSections
	lines := util.StringSplitLines(s)
	for _, line := range lines {
		if line != "" {
			switch line[0] {
			case ' ':
				ret = ret.AddLine(patchTypeUnchanged, line)
			case '+':
				ret = ret.AddLine(patchTypeAdded, line)
			case '-':
				ret = ret.AddLine(patchTypeRemoved, line)
			default:
				ret = ret.AddLine("location", line)
			}
		}
	}
	return ret
}
