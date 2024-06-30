package cproject

import (
	"fmt"
	"net/http"

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
	depsKey = "deps"
	pkgsKey = "packages"
)

func RunAction(w http.ResponseWriter, r *http.Request) {
	actQ, _ := cutil.PathString(r, "act", false)
	act := "run.action." + actQ
	if phase := r.URL.Query().Get("phase"); phase != "" {
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
		if actT.Expensive(cfg) {
			if cfg.GetStringOpt("hasloaded") != util.BoolTrue {
				cutil.URLAddQuery(r.URL, "hasloaded", util.BoolTrue)
				page := &vpage.Load{URL: r.URL.String(), Title: fmt.Sprintf("Running [%s] for [%s]", phase, prj.Title())}
				return controller.Render(r, as, page, ps, "projects", prj.Key, actT.Breadcrumb())
			}
		}
		if isBuild && phase == "" {
			ps.Data = action.AllBuilds
			page := &vbuild.BuildResult{Project: prj, Cfg: cfg}
			return controller.Render(r, as, page, ps, "projects", prj.Key, actT.Breadcrumb())
		}
		result := action.Apply(ps.Context, actionParams(prj.Key, actT, cfg, as, ps.Logger))
		if result == nil {
			result = &action.Result{}
		}
		if result.Project == nil {
			result.Project = prj
		}
		if redir, ok := result.Data.(string); ok {
			return redir, nil
		}
		ps.Data = result
		if isBuild {
			if phase == depsKey {
				return runDeps(prj, result, w, r, as, ps)
			}
			if phase == pkgsKey {
				return runPkgs(prj, result, w, r, as, ps)
			}
			page := &vbuild.BuildResult{Project: prj, Cfg: cfg, BuildResult: result}
			return controller.Render(r, as, page, ps, "projects", prj.Key, actT.Breadcrumb())
		}

		page := &vaction.Result{Ctx: &action.ResultContext{Prj: prj, Cfg: cfg, Res: result}}
		return controller.Render(r, as, page, ps, "projects", prj.Key, actT.Breadcrumb())
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

	cfg := cutil.QueryArgsMap(r)
	prj, err := as.Services.Projects.Get(tgt)
	if err != nil {
		return nil, actT, nil, err
	}
	return cfg, actT, prj, nil
}
