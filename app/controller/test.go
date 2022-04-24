package controller

import (
	"github.com/pkg/errors"
	"github.com/valyala/fasthttp"

	"projectforge.dev/projectforge/app"
	"projectforge.dev/projectforge/app/action"
	"projectforge.dev/projectforge/app/controller/cutil"
	"projectforge.dev/projectforge/app/diff"
	"projectforge.dev/projectforge/app/util"
	"projectforge.dev/projectforge/views/layout"
	"projectforge.dev/projectforge/views/vaction"
	"projectforge.dev/projectforge/views/vtest"
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
			res := action.Apply(ps.Context, actionParams("testproject", action.TypeTest, cfg, as, ps.Logger))
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
