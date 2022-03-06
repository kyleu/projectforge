package controller

import (
	"github.com/pkg/errors"
	"github.com/valyala/fasthttp"
	"projectforge.dev/app/action"
	"projectforge.dev/app/diff"
	"projectforge.dev/app/lib/telemetry"
	"projectforge.dev/app/util"
	"projectforge.dev/views/layout"
	"projectforge.dev/views/vaction"
	"projectforge.dev/views/vtest"

	"projectforge.dev/app/controller/cutil"

	"projectforge.dev/app"
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
		key, err := RCRequiredString(rc, "key", false)
		if err != nil {
			return "", err
		}
		ps.Title = "Test [" + key + "]"
		ps.Data = key

		var page layout.Page
		switch key {
		case "diff":
			ret := diff.Results{}
			for _, x := range diff.AllExamples {
				res := x.Calc()
				ret = append(ret, res)
			}
			ps.Data = ret
			page = &vtest.Diffs{Results: ret}
		case "bootstrap":
			cfg := util.ValueMap{}
			cfg.Add("path", "./testproject", "method", key, "wipe", true)
			nc, span, logger := telemetry.StartSpan(ps.Context, "action:test.run", ps.Logger)
			res := action.Apply(nc, actionParams(span, "testproject", action.TypeTest, cfg, as, logger))
			ps.Data = res

			_, err = as.Services.Projects.Refresh()
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
