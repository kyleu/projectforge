package ctest

import (
	"net/http"

	"projectforge.dev/projectforge/app/controller/cutil"
	"projectforge.dev/projectforge/app/doctor/libraries"
	"projectforge.dev/projectforge/app/util"
	"projectforge.dev/projectforge/views/layout"
	"projectforge.dev/projectforge/views/vtest"
)

func librariesTest(r *http.Request, ps *cutil.PageState) (layout.Page, error) {
	var ret libraries.Results
	act := util.OrDefault(cutil.QueryStringString(ps.URI, "act"), "test")
	switch tgt := cutil.QueryStringString(ps.URI, "tgt"); tgt {
	case "":
		ps.SetTitleAndData("Library Tests", libraries.AllLibraries)
	case "all":
		var errs []error
		ret, errs = util.AsyncCollect(libraries.AllLibraries, func(lib *libraries.Library) (*libraries.Result, error) {
			return libraries.Process(ps.Context, lib, act, ps.Logger)
		})
		if len(errs) != 0 {
			return nil, util.ErrorMerge(errs...)
		}
		ps.SetTitleAndData("All Library Tests", ret)
	default:
		lib := libraries.AllLibraries.Get(tgt)
		x, err := libraries.Process(ps.Context, lib, act, ps.Logger)
		if err != nil {
			return nil, err
		}
		ret = append(ret, x)
	}
	ret.Sort()
	page := &vtest.LibraryResult{Results: ret}
	return page, nil
}
