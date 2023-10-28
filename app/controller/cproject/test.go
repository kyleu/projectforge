package cproject

import (
	"github.com/pkg/errors"
	"github.com/samber/lo"
	"github.com/valyala/fasthttp"

	"projectforge.dev/projectforge/app"
	"projectforge.dev/projectforge/app/controller"
	"projectforge.dev/projectforge/app/controller/cutil"
	"projectforge.dev/projectforge/app/file/diff"
	"projectforge.dev/projectforge/app/project/action"
	"projectforge.dev/projectforge/app/util"
	"projectforge.dev/projectforge/views/layout"
	"projectforge.dev/projectforge/views/vaction"
	"projectforge.dev/projectforge/views/vtest"
)

func TestList(rc *fasthttp.RequestCtx) {
	controller.Act("test.list", rc, func(as *app.State, ps *cutil.PageState) (string, error) {
		ps.SetTitleAndData("Tests", []string{"bootstrap", "diff"})
		return controller.Render(rc, as, &vtest.List{}, ps, "Tests")
	})
}

func TestRun(rc *fasthttp.RequestCtx) {
	controller.Act("test.run", rc, func(as *app.State, ps *cutil.PageState) (string, error) {
		key, err := cutil.RCRequiredString(rc, "key", false)
		if err != nil {
			return "", err
		}
		ps.SetTitleAndData("Test ["+key+"]", key)
		bc := []string{"Tests||/test"}

		var page layout.Page
		switch key {
		case "diff":
			ret := lo.Map(diff.AllExamples, func(x *diff.Example, _ int) *diff.Result {
				return x.Calc()
			})
			bc = append(bc, "Diff")
			ps.SetTitleAndData("Diff Test", ret)
			page = &vtest.Diffs{Results: ret}
		case "bootstrap":
			cfg := util.ValueMap{}
			cfg.Add("path", "./testproject", "method", key, "wipe", true)
			res := action.Apply(ps.Context, actionParams("testproject", action.TypeTest, cfg, as, ps.Logger))

			bc = append(bc, "Bootstrap")
			ps.Data = res

			_, err = as.Services.Projects.Refresh(ps.Logger)
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
		return controller.Render(rc, as, page, ps, bc...)
	})
}
