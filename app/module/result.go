package module

import (
	"projectforge.dev/projectforge/app/file/diff"
	"projectforge.dev/projectforge/app/util"
)

type Result struct {
	Keys     []string    `json:"keys"`
	Status   string      `json:"status"`
	Diffs    diff.Diffs  `json:"diffs,omitempty"`
	Actions  Resolutions `json:"actions,omitempty"`
	Duration int         `json:"duration,omitempty"`
}

func (r *Result) DiffsFiltered(includeSkipped bool) diff.Diffs {
	ret := make(diff.Diffs, 0, len(r.Diffs))
	for _, d := range r.Diffs {
		if includeSkipped || d.Status != diff.StatusSkipped {
			ret = append(ret, d)
		}
	}
	return ret
}

type Results []*Result

func (r Results) DiffCount(includeSkipped bool) int {
	ret := 0
	for _, m := range r {
		ret += len(m.DiffsFiltered(includeSkipped))
	}
	return ret
}

func (r Results) Paths(includeSkipped bool) []string {
	ret := make([]string, 0, r.DiffCount(includeSkipped))
	for _, res := range r {
		ret = append(ret, res.DiffsFiltered(includeSkipped).Paths()...)
	}
	return util.ArrayRemoveDuplicates(ret)
}
