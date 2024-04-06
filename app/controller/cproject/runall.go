package cproject

import (
	"fmt"
	"net/http"

	"github.com/samber/lo"

	"projectforge.dev/projectforge/app"
	"projectforge.dev/projectforge/app/controller"
	"projectforge.dev/projectforge/app/controller/cutil"
	"projectforge.dev/projectforge/app/project"
	"projectforge.dev/projectforge/app/project/action"
	"projectforge.dev/projectforge/app/util"
	"projectforge.dev/projectforge/views"
	"projectforge.dev/projectforge/views/vaction"
)

func RunAllActions(w http.ResponseWriter, r *http.Request) {
	helpKey := "run.all"
	actKey, _ := cutil.PathString(r, "act", false)
	if actKey != "" {
		helpKey += "." + actKey
	}
	controller.Act(helpKey, w, r, func(as *app.State, ps *cutil.PageState) (string, error) {
		actS, err := cutil.PathString(r, "act", false)
		if err != nil {
			return "", err
		}
		cfg := cutil.QueryArgsMap(r)
		prjs := as.Services.Projects.Projects()
		tags := util.StringSplitAndTrim(r.URL.Query().Get("tags"), ",")
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
			return runAllStart(w, r, as, ps)
		}
		actT := action.TypeFromString(actS)

		if actT.Key == action.TypeBuild.Key {
			switch cfg.GetStringOpt("phase") {
			case "":
				ps.SetTitleAndData("Build All Projects", prjs)
				page := &vaction.Results{T: actT, Cfg: cfg, Projects: prjs, Ctxs: nil, Tags: tags, IsBuild: true}
				return controller.Render(w, r, as, page, ps, "projects", actT.Title)
			case depsKey:
				return runAllDeps(cfg, prjs, tags, w, r, as, ps)
			case pkgsKey:
				return runAllPkgs(cfg, prjs, w, r, as, ps)
			}
		}

		results := action.ApplyAll(ps.Context, prjs, actT, cfg, as, ps.Logger)

		ps.SetTitleAndData(fmt.Sprintf("[%s] All Projects", actT.Title), results)
		page := &vaction.Results{T: actT, Cfg: cfg, Projects: prjs, Ctxs: results, Tags: tags}
		return controller.Render(w, r, as, page, ps, "projects", actT.Title)
	})
}

func runAllStart(w http.ResponseWriter, r *http.Request, as *app.State, ps *cutil.PageState) (string, error) {
	ps.SetTitleAndData("Start All", "TODO")
	return controller.Render(w, r, as, &views.Debug{}, ps, "projects", "Start")
}
