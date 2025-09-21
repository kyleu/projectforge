package module

import (
	"github.com/samber/lo"

	"projectforge.dev/projectforge/app/file/diff"
)

type Result struct {
	Keys     []string    `json:"keys"`
	Status   string      `json:"status"`
	Diffs    diff.Diffs  `json:"diffs,omitempty"`
	Actions  Resolutions `json:"actions,omitempty"`
	Duration int         `json:"duration,omitzero"`
}

func (r *Result) DiffsFiltered(includeSkipped bool) diff.Diffs {
	return lo.Filter(r.Diffs, func(d *diff.Diff, _ int) bool {
		return includeSkipped || d.Status != diff.StatusSkipped
	})
}

type Results []*Result

func (r Results) DiffCount(includeSkipped bool) int {
	return lo.Sum(lo.Map(r, func(m *Result, _ int) int {
		return len(m.DiffsFiltered(includeSkipped))
	}))
}

func (r Results) Paths(includeSkipped bool) []string {
	return lo.Uniq(lo.FlatMap(r, func(res *Result, _ int) []string {
		return res.DiffsFiltered(includeSkipped).Paths()
	}))
}
