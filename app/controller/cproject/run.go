package cproject

import (
	"fmt"

	"github.com/valyala/fasthttp"

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

func RunAction(rc *fasthttp.RequestCtx) {
	actQ, _ := cutil.RCRequiredString(rc, "act", false)
	act := "run.action." + actQ
	if phase := string(rc.URI().QueryArgs().Peek("phase")); phase != "" {
		act += "." + phase
	}
	controller.Act(act, rc, func(as *app.State, ps *cutil.PageState) (string, error) {
		cfg, actT, prj, err := loadActionProject(rc, as)
		if err != nil {
			return "", err
		}
		ps.Title = fmt.Sprintf("[%s] %s", actT.Title, prj.Title())
		ps.Data = actT
		if curr, ok := cfg["path"]; !ok || curr == "" {
			cfg["path"] = prj.Path
		}
		isBuild := actT.Key == action.TypeBuild.Key
		phase := cfg.GetStringOpt("phase")
		if actT.Expensive(cfg) {
			if cfg.GetStringOpt("hasloaded") != util.BoolTrue {
				rc.URI().QueryArgs().Set("hasloaded", util.BoolTrue)
				page := &vpage.Load{URL: rc.URI().String(), Title: fmt.Sprintf("Running [%s] for [%s]", phase, prj.Title())}
				return controller.Render(rc, as, page, ps, "projects", prj.Key, actT.Title)
			}
		}
		if isBuild && phase == "" {
			ps.Data = action.AllBuilds
			page := &vbuild.BuildResult{Project: prj, Cfg: cfg, BuildResult: nil, GitResult: nil}
			return controller.Render(rc, as, page, ps, "projects", prj.Key, actT.Title)
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
				return runDeps(prj, result, rc, as, ps)
			}
			if phase == pkgsKey {
				return runPkgs(prj, result, rc, as, ps)
			}
			page := &vbuild.BuildResult{Project: prj, Cfg: cfg, BuildResult: result}
			return controller.Render(rc, as, page, ps, "projects", prj.Key, actT.Title)
		}

		page := &vaction.Result{Ctx: &action.ResultContext{Prj: prj, Cfg: cfg, Res: result}}
		return controller.Render(rc, as, page, ps, "projects", prj.Key, actT.Title)
	})
}

func loadActionProject(rc *fasthttp.RequestCtx, as *app.State) (util.ValueMap, action.Type, *project.Project, error) {
	actS, err := cutil.RCRequiredString(rc, "act", false)
	if err != nil {
		return nil, action.TypeTest, nil, err
	}
	actT := action.TypeFromString(actS)
	tgt, err := cutil.RCRequiredString(rc, "key", false)
	if err != nil {
		return nil, actT, nil, err
	}

	cfg := cutil.QueryArgsMap(rc)
	prj, err := as.Services.Projects.Get(tgt)
	if err != nil {
		return nil, actT, nil, err
	}
	return cfg, actT, prj, nil
}
