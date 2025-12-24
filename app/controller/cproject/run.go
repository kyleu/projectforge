package cproject

import (
	"fmt"
	"net/http"

	"github.com/pkg/errors"

	"projectforge.dev/projectforge/app"
	"projectforge.dev/projectforge/app/controller"
	"projectforge.dev/projectforge/app/controller/cutil"
	"projectforge.dev/projectforge/app/project"
	"projectforge.dev/projectforge/app/project/action"
	"projectforge.dev/projectforge/app/util"
	"projectforge.dev/projectforge/views/vaction"
	"projectforge.dev/projectforge/views/vbuild"
	"projectforge.dev/projectforge/views/vpage"
)

const (
	depsKey     = "deps"
	pkgsKey     = "packages"
	statsKey    = "codestats"
	coverageKey = "coverage"
	keyCustom   = "custom"
)

//nolint:gocognit,nestif
func RunAction(w http.ResponseWriter, r *http.Request) {
	actQ, _ := cutil.PathString(r, "act", false)
	act := "run.action." + actQ
	if phase := cutil.QueryStringString(r, "phase"); phase != "" {
		act += "." + phase
	}
	controller.Act(act, w, r, func(as *app.State, ps *cutil.PageState) (string, error) {
		cfg, actT, prj, err := loadActionProject(r, as)
		if err != nil {
			return "", err
		}
		ps.SetTitleAndData(fmt.Sprintf("[%s] %s", actT.Title, prj.Title()), actT)
		if curr, ok := cfg["path"]; !ok || curr == "" {
			cfg["path"] = prj.Path
		}
		isBuild := actT.Matches(action.TypeBuild)
		phase := cfg.GetStringOpt("phase")
		bc := []string{"projects", prj.Key, actT.Breadcrumb()}
		if actT.Expensive(cfg) {
			if page := HandleLoad(cfg, r.URL, fmt.Sprintf("Running [%s] for [%s]", actT.Title, prj.Title())); page != nil {
				return controller.Render(r, as, page, ps, bc...)
			}
		}
		if isBuild {
			if phase == "" {
				ps.SetTitleAndData("Build Result", action.AllBuilds)
				page := &vbuild.BuildResult{Project: prj, Cfg: cfg}
				return controller.Render(r, as, page, ps, bc...)
			}
			b := action.AllBuilds.Get(phase)
			if b == nil {
				return "", errors.Errorf("unknown phase [%s]", phase)
			}
			if b.Key == keyCustom {
				if cmd := cfg.GetStringOpt("cmd"); cmd == "" {
					argRes := util.FieldDescsCollectMap(cfg, gitCustomArgs)
					if argRes.HasMissing() {
						ps.SetTitleAndData("Custom Command", argRes)
						url := fmt.Sprintf("/run/%s/build", prj.Key)
						page := &vpage.Args{URL: url, Directions: "Enter your commit message", Results: argRes, Hidden: map[string]string{"phase": phase}}
						return controller.Render(r, as, page, ps, bc...)
					}
				}
			}
			if b.Expensive {
				if page := HandleLoad(cfg, r.URL, fmt.Sprintf("Running [%s] for [%s]", phase, prj.Title())); page != nil {
					return controller.Render(r, as, page, ps, bc...)
				}
			}
		}
		result := action.Apply(ps.Context, actionParams(prj.Key, actT, cfg, as, ps.Logger))
		if result == nil {
			result = &action.Result{}
		}
		if result.Project == nil {
			result.Project = prj
		}
		if redir, e := util.Cast[string](result.Data); e == nil {
			return redir, nil
		}
		ps.SetTitleAndData("Action ["+actT.String()+"]", result)
		if isBuild {
			if phase == depsKey {
				return runDeps(prj, result, r, as, ps)
			}
			if phase == pkgsKey {
				return runPkgs(prj, result, r, as, ps)
			}
			if phase == statsKey {
				return runCodeStats(prj, result, r, as, ps)
			}
			if phase == coverageKey {
				return runCoverage(prj, result, r, as, ps)
			}
			page := &vbuild.BuildResult{Project: prj, Cfg: cfg, BuildResult: result}
			return controller.Render(r, as, page, ps, bc...)
		}

		page := &vaction.Result{Ctx: &action.ResultContext{Prj: prj, Cfg: cfg, Res: result}}
		return controller.Render(r, as, page, ps, bc...)
	})
}

func loadActionProject(r *http.Request, as *app.State) (util.ValueMap, action.Type, *project.Project, error) {
	actS, err := cutil.PathString(r, "act", false)
	if err != nil {
		return nil, action.TypeTest, nil, err
	}
	actT := action.TypeFromString(actS)
	tgt, err := cutil.PathString(r, "key", false)
	if err != nil {
		return nil, actT, nil, err
	}

	cfg := cutil.QueryStringAsMap(r.URL)
	prj, err := as.Services.Projects.Get(tgt)
	if err != nil {
		return nil, actT, nil, err
	}
	return cfg, actT, prj, nil
}
