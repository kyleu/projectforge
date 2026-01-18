package cproject

import (
	"fmt"
	"net/http"

	"github.com/pkg/errors"
	"github.com/samber/lo"

	"projectforge.dev/projectforge/app"
	"projectforge.dev/projectforge/app/controller"
	"projectforge.dev/projectforge/app/controller/cutil"
	"projectforge.dev/projectforge/app/project"
	"projectforge.dev/projectforge/app/project/action"
	"projectforge.dev/projectforge/app/util"
	"projectforge.dev/projectforge/views"
	"projectforge.dev/projectforge/views/vaction"
	"projectforge.dev/projectforge/views/vpage"
)

//nolint:gocognit
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
		cfg := cutil.QueryStringAsMap(r.URL)
		prjs := as.Services.Projects.Projects()
		tags := util.StringSplitAndTrim(cutil.QueryStringString(ps.URI, "tags"), ",")
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

		if actT.Expensive(cfg) {
			x := actT.String()
			if actT.Matches(action.TypeBuild) {
				x += ":" + cfg.GetStringOpt("phase")
			}
			if page := HandleLoad(cfg, r.URL, fmt.Sprintf("Running [%s] action for all projects", x)); page != nil {
				return controller.Render(r, as, page, ps, "projects", actT.Breadcrumb())
			}
		}

		if actT.Matches(action.TypeBuild) {
			phase := cfg.GetStringOpt("phase")
			if phase == keyCustom {
				if cmd := cfg.GetStringOpt("cmd"); cmd == "" {
					argRes := util.FieldDescsCollectMap(cfg, gitCustomArgs)
					if argRes.HasMissing() {
						ps.SetTitleAndData("Custom Command", argRes)
						page := &vpage.Args{URL: "/run/build", Directions: "Enter your commit message", Results: argRes, Hidden: map[string]string{"phase": phase}}
						return controller.Render(r, as, page, ps, actT.Breadcrumb())
					}
				}
			}

			switch phase {
			case "":
				ps.SetTitleAndData("Build All Projects", prjs)
				page := &vaction.Results{T: actT, Cfg: cfg, Projects: prjs, Ctxs: nil, Tags: tags, IsBuild: true}
				return controller.Render(r, as, page, ps, "projects", actT.Breadcrumb())
			case depsKey:
				return runAllDeps(cfg, prjs, tags, r, as, ps)
			case pkgsKey:
				return runAllPkgs(cfg, prjs, r, as, ps)
			case statsKey:
				return runAllCodeStats(cfg, prjs, r, as, ps)
			case coverageKey:
				return runAllCoverage(cfg, prjs, r, as, ps)
			case "lint":
				return "", errors.New("can't run multiple instances of golangcilint")
			}
		}

		t := util.TimerStart()
		results := action.ApplyAll(ps.Context, prjs, actT, cfg, as, ps.Logger)

		ps.SetTitleAndData(fmt.Sprintf("[%s] All Projects", actT.Title), results)
		page := &vaction.Results{T: actT, Cfg: cfg, Projects: prjs, Ctxs: results, Tags: tags, Duration: t.EndString()}
		return controller.Render(r, as, page, ps, "projects", actT.Breadcrumb())
	})
}

func runAllStart(w http.ResponseWriter, r *http.Request, as *app.State, ps *cutil.PageState) (string, error) {
	ps.SetTitleAndData("Start All", "Starting all projects at once is not available yet. Use the Start action on each project.")
	return controller.Render(r, as, &views.Debug{}, ps, "projects", "Start**play")
}
