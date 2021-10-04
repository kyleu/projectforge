package controller

import (
	"github.com/kyleu/projectforge/app/action"
	"github.com/kyleu/projectforge/app/diff"
	"github.com/kyleu/projectforge/app/telemetry"
	"github.com/kyleu/projectforge/app/util"
	"github.com/kyleu/projectforge/views/layout"
	"github.com/kyleu/projectforge/views/vaction"
	"github.com/kyleu/projectforge/views/vtest"
	"github.com/pkg/errors"
	"github.com/valyala/fasthttp"

	"github.com/kyleu/projectforge/app/controller/cutil"

	"github.com/kyleu/projectforge/app"
)

func TestList(rc *fasthttp.RequestCtx) {
	act("test.list", rc, func(as *app.State, ps *cutil.PageState) (string, error) {
		ps.Title = "Tests"
		ps.Data = "TODO"
		return render(rc, as, &vtest.List{}, ps, "Tests")
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

		var page layout.Page
		switch key {
		case "diff":
			ret := []*diff.Result{}
			for _, x := range diff.AllExamples {
				res := x.Calc()
				ret = append(ret, res)
			}
			ps.Data = ret
			page = &vtest.Diffs{Results: ret}
		case "bootstrap":
			cfg := util.ValueMap{}
			cfg.Add("path", "./testproject", "method", key, "wipe", true)
			nc, span := telemetry.StartSpan(ps.Context, "action", "test.run")
			res := action.Apply(nc, &action.Params{
				Span: span, ProjectKey: "testproject", T: action.TypeTest, Cfg: cfg, RootFiles: as.Files,
				MSvc: as.Services.Modules, PSvc: as.Services.Projects, Logger: ps.Logger,
			})
			ps.Data = res

			_, err = as.Services.Projects.Refresh(as.Files)
			if err != nil {
				return "", err
			}

			prj, err := as.Services.Projects.Get("testproject")
			if err != nil {
				return "", err
			}

			page = &vaction.Result{Ctx: &action.ResultContext{Prj: prj, Cfg: cfg, Res: res}}
		default:
			return "", errors.New("invalid test [" + key + "]")
		}
		return render(rc, as, page, ps, "Tests", key)
	})
}
