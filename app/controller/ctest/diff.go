package ctest

import (
	"github.com/samber/lo"

	"projectforge.dev/projectforge/app/controller/cutil"
	"projectforge.dev/projectforge/app/file/diff"
	"projectforge.dev/projectforge/views/layout"
	"projectforge.dev/projectforge/views/vtest"
)

func diffTest(ps *cutil.PageState) (layout.Page, error) {
	ret := lo.Map(diff.AllExamples, func(x *diff.Example, _ int) *diff.Result {
		return x.Calc()
	})
	ps.SetTitleAndData("Diff Test", ret)
	return &vtest.Diffs{Results: ret}, nil
}
