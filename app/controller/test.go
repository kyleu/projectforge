package controller

import (
	"github.com/kyleu/projectforge/app/action"
	"github.com/kyleu/projectforge/app/telemetry"
	"github.com/kyleu/projectforge/app/util"
	"github.com/kyleu/projectforge/views"
	"github.com/kyleu/projectforge/views/vaction"
	"github.com/valyala/fasthttp"

	"github.com/kyleu/projectforge/app/controller/cutil"

	"github.com/kyleu/projectforge/app"
)

func TestList(rc *fasthttp.RequestCtx) {
	act("test.list", rc, func(as *app.State, ps *cutil.PageState) (string, error) {
		ps.Title = "Tests"
		ps.Data = "TODO"
		return render(rc, as, &views.Debug{}, ps, "Search")
	})
}

func TestRun(rc *fasthttp.RequestCtx) {
	act("test.run", rc, func(as *app.State, ps *cutil.PageState) (string, error) {
		key, err := rcRequiredString(rc, "key", false)
		if err != nil {
			return "", err
		}
		ps.Title = "Test [" + key + "]"
		ps.Data = key
		cfg := util.ValueMap{}
		cfg.Add("path", "./testproject", "method", key, "wipe", true)
		nc, span := telemetry.StartSpan(ps.Context, "action", "test.run")
		res := action.Apply(nc, span, "testproject", action.TypeTest, cfg, as.Services.Modules, as.Services.Projects, ps.Logger)
		ps.Data = res

		_, err = as.Services.Projects.Refresh()
		if err != nil {
			return "", err
		}

		prj, err := as.Services.Projects.Get("testproject")
		if err != nil {
			return "", err
		}

		page := &vaction.Result{Ctx: &action.ResultContext{Prj: prj, Cfg: cfg, Res: res}}
		return render(rc, as, page, ps, "Bootstrap")
	})
}
