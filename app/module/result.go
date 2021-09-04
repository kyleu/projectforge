package module

import (
	"github.com/kyleu/projectforge/app/diff"
)

type Result struct {
	Keys     []string     `json:"keys"`
	Status   string       `json:"status"`
	Diffs    []*diff.Diff `json:"diffs,omitempty"`
	Actions  Resolutions  `json:"actions,omitempty"`
	Duration int          `json:"duration,omitempty"`
}

func (r *Result) DiffsFiltered(includeSkipped bool) []*diff.Diff {
	ret := make([]*diff.Diff, 0, len(r.Diffs))
	for _, d := range r.Diffs {
		if includeSkipped || d.Status != diff.StatusSkipped {
			ret = append(ret, d)
		}
	}
	return ret
}

type Results []*Result
