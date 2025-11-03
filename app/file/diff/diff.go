package diff

import (
	"fmt"

	"github.com/samber/lo"

	"projectforge.dev/projectforge/app/file"
	"projectforge.dev/projectforge/app/util"
)

type Diff struct {
	Path    string  `json:"path"`
	Status  *Status `json:"status"`
	Patch   string  `json:"patch,omitzero"`
	Changes Changes `json:"changes,omitempty"`
}

func (d *Diff) String() string {
	return fmt.Sprintf("%s:%s", d.Path, d.Status)
}

type Diffs []*Diff

func (d Diffs) HasStatus(s *Status) bool {
	return lo.ContainsBy(d, func(x *Diff) bool {
		return x.Status.Matches(s)
	})
}

func (d Diffs) Paths() []string {
	return lo.Uniq(lo.Map(d, func(x *Diff, _ int) string {
		return x.Path
	}))
}

func File(src *file.File, tgt *file.File) *Diff {
	ret := &Diff{Path: src.FullPath()}
	switch {
	case src == nil:
		ret.Status = StatusMissing
	case tgt == nil:
		ret.Status = StatusNew
		ret.Patch = fmt.Sprintf("[new file, %s]", util.ByteSizeSI(int64(len(src.Content))))
	default:
		d := Calc(ret.Path, src.Content, tgt.Content)
		if len(d.Changes) > 0 {
			ret.Patch = d.Patch
			ret.Status = StatusDifferent
			ret.Changes = d.Changes
		} else {
			ret.Status = StatusIdentical
		}
	}
	return ret
}

func Files(src file.Files, tgt file.Files, includeUnchanged bool) Diffs {
	return lo.FlatMap(src, func(s *file.File, _ int) []*Diff {
		p := s.FullPath()
		t := tgt.Get(p)
		d := File(s, t)
		if includeUnchanged || d.Status != StatusIdentical {
			return Diffs{d}
		}
		return nil
	})
}
