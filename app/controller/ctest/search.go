package ctest

import (
	"net/http"

	"github.com/samber/lo"

	"projectforge.dev/projectforge/app"
	"projectforge.dev/projectforge/app/controller/cutil"
	"projectforge.dev/projectforge/app/project"
	"projectforge.dev/projectforge/app/project/action"
	"projectforge.dev/projectforge/app/util"
	"projectforge.dev/projectforge/views/layout"
	"projectforge.dev/projectforge/views/vaction"
)

func searchTest(as *app.State, r *http.Request, ps *cutil.PageState) (layout.Page, error) {
	q := r.URL.Query().Get("q")
	cfg := util.ValueMap{"q": q}

	prjs := as.Services.Projects.Projects().WithModules("export")
	ctxs := lo.Map(prjs, func(p *project.Project, _ int) *action.ResultContext {
		ret := util.OK
		res := &action.Result{Project: p, Action: action.TypeTest, Status: "ok", Args: cfg, Data: ret}
		return &action.ResultContext{Prj: p, Res: res}
	})

	ps.SetTitleAndData("Search", q)
	page := &vaction.Results{Projects: prjs, T: action.TypeTest, Cfg: cfg, Ctxs: ctxs}
	return page, nil
}
