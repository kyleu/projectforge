package cproject

import (
	"net/http"

	"github.com/pkg/errors"
	"github.com/samber/lo"

	"projectforge.dev/projectforge/app"
	"projectforge.dev/projectforge/app/controller"
	"projectforge.dev/projectforge/app/controller/cutil"
	"projectforge.dev/projectforge/app/file/diff"
	"projectforge.dev/projectforge/app/project"
	"projectforge.dev/projectforge/app/project/action"
	"projectforge.dev/projectforge/app/util"
	"projectforge.dev/projectforge/views/layout"
	"projectforge.dev/projectforge/views/vaction"
	"projectforge.dev/projectforge/views/vtest"
)

func TestList(w http.ResponseWriter, r *http.Request) {
	controller.Act("test.list", w, r, func(as *app.State, ps *cutil.PageState) (string, error) {
		ps.SetTitleAndData("Tests", []string{"bootstrap", "diff"})
		return controller.Render(r, as, &vtest.List{}, ps, "Tests")
	})
}

func TestRun(w http.ResponseWriter, r *http.Request) {
	controller.Act("test.run", w, r, func(as *app.State, ps *cutil.PageState) (string, error) {
		key, err := cutil.PathString(r, "key", false)
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
			ps.SetTitleAndData("Bootstrap", res)

			_, err = as.Services.Projects.Refresh(ps.Logger)
			if err != nil {
				return "", err
			}

			prj, err := as.Services.Projects.Get("testproject")
			if err != nil {
				return "", err
			}

			page = &vaction.Result{Ctx: &action.ResultContext{Prj: prj, Cfg: cfg, Res: res}}
		case "search":
			q := r.URL.Query().Get("q")
			cfg := util.ValueMap{"q": q}

			prjs := as.Services.Projects.Projects().WithModules("export")
			ctxs := lo.Map(prjs, func(p *project.Project, _ int) *action.ResultContext {
				ret := util.OK
				res := &action.Result{Project: p, Action: action.TypeTest, Status: "ok", Args: cfg, Data: ret}
				return &action.ResultContext{Prj: p, Res: res}
			})

			bc = append(bc, "Search")
			ps.SetTitleAndData("Search", q)
			page = &vaction.Results{Projects: prjs, T: action.TypeTest, Cfg: cfg, Ctxs: ctxs}
		default:
			return "", errors.Errorf("invalid test [%s]", key)
		}
		return controller.Render(r, as, page, ps, bc...)
	})
}
