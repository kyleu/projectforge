package cproject

import (
	"fmt"

	"github.com/samber/lo"
	"github.com/valyala/fasthttp"

	"projectforge.dev/projectforge/app"
	"projectforge.dev/projectforge/app/controller"
	"projectforge.dev/projectforge/app/controller/cutil"
	"projectforge.dev/projectforge/app/project"
	"projectforge.dev/projectforge/app/project/action"
	"projectforge.dev/projectforge/app/util"
	"projectforge.dev/projectforge/views"
	"projectforge.dev/projectforge/views/vaction"
)

func RunAllActions(rc *fasthttp.RequestCtx) {
	helpKey := "run.all"
	actKey, _ := cutil.RCRequiredString(rc, "act", false)
	if actKey != "" {
		helpKey += "." + actKey
	}
	controller.Act(helpKey, rc, func(as *app.State, ps *cutil.PageState) (string, error) {
		actS, err := cutil.RCRequiredString(rc, "act", false)
		if err != nil {
			return "", err
		}
		cfg := cutil.QueryArgsMap(rc)
		prjs := as.Services.Projects.Projects()
		tags := util.StringSplitAndTrim(string(rc.URI().QueryArgs().Peek("tags")), ",")
		if len(tags) == 0 {
			prjs = prjs.WithoutTags("all-skip")
		} else {
			filtered := prjs.WithTags(tags...)
			if len(filtered) == 0 && len(tags) == 1 {
				key := tags[0]
				if key[0] == '-' {
					prjs = lo.Filter(prjs, func(x *project.Project, _ int) bool {
						return x.Key != key[1:]
					})
				} else {
					prjs = project.Projects{prjs.Get(key)}
				}
			} else {
				prjs = filtered
			}
		}
		if actS == "start" {
			return runAllStart(rc, as, ps)
		}
		actT := action.TypeFromString(actS)

		if actT.Key == action.TypeBuild.Key {
			switch cfg.GetStringOpt("phase") {
			case "":
				ps.SetTitleAndData("Build All Projects", prjs)
				page := &vaction.Results{T: actT, Cfg: cfg, Projects: prjs, Ctxs: nil, Tags: tags, IsBuild: true}
				return controller.Render(rc, as, page, ps, "projects", actT.Title)
			case depsKey:
				return runAllDeps(cfg, prjs, tags, rc, as, ps)
			case pkgsKey:
				return runAllPkgs(cfg, prjs, rc, as, ps)
			}
		}

		results := action.ApplyAll(ps.Context, prjs, actT, cfg, as, ps.Logger)

		ps.SetTitleAndData(fmt.Sprintf("[%s] All Projects", actT.Title), results)
		page := &vaction.Results{T: actT, Cfg: cfg, Projects: prjs, Ctxs: results, Tags: tags}
		return controller.Render(rc, as, page, ps, "projects", actT.Title)
	})
}

func runAllStart(rc *fasthttp.RequestCtx, as *app.State, ps *cutil.PageState) (string, error) {
	ps.SetTitleAndData("Start All", "TODO")
	return controller.Render(rc, as, &views.Debug{}, ps, "projects", "Start")
}
