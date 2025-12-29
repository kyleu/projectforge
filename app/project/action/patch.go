package action

import (
	"github.com/samber/lo"

	"projectforge.dev/projectforge/app/util"
)

type PatchSection struct {
	Type  string   `json:"type"`
	Lines []string `json:"lines"`
}

func (p *PatchSection) CSSClass() string {
	switch p.Type {
	case "unchanged":
		return "color-muted"
	case "added":
		return "success"
	case "removed":
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
	return p.Type == "added" || p.Type == "removed"
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
				ret = ret.AddLine("unchanged", line)
			case '+':
				ret = ret.AddLine("added", line)
			case '-':
				ret = ret.AddLine("removed", line)
			default:
				ret = ret.AddLine("location", line)
			}
		}
	}
	return ret
}
