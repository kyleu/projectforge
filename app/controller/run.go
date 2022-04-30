package controller

import (
	"fmt"
	"strings"

	"github.com/valyala/fasthttp"
	"golang.org/x/exp/slices"
	"projectforge.dev/projectforge/app"
	"projectforge.dev/projectforge/app/action"
	"projectforge.dev/projectforge/app/controller/cutil"
	"projectforge.dev/projectforge/app/project"
	"projectforge.dev/projectforge/app/util"
	"projectforge.dev/projectforge/views/vaction"
	"projectforge.dev/projectforge/views/vbuild"
)

func RunAction(rc *fasthttp.RequestCtx) {
	act("run.action", rc, func(as *app.State, ps *cutil.PageState) (string, error) {
		tgt, err := RCRequiredString(rc, "key", false)
		if err != nil {
			return "", err
		}
		actS, err := RCRequiredString(rc, "act", false)
		if err != nil {
			return "", err
		}

		cfg := util.ValueMap{}
		actT := action.TypeFromString(actS)
		prj, err := as.Services.Projects.Get(tgt)
		if err != nil {
			return "", err
		}
		cfg["path"] = prj.Path
		rc.QueryArgs().VisitAll(func(k []byte, v []byte) {
			cfg[string(k)] = string(v)
		})

		isBuild := actT.Key == action.TypeBuild.Key
		phase := cfg.GetStringOpt("phase")

		if isBuild && phase == "" {
			page := &vbuild.BuildResult{Project: prj, BuildResult: nil, GitResult: nil}
			return render(rc, as, page, ps, "projects", actT.Title)
		}

		result := action.Apply(ps.Context, actionParams(tgt, actT, cfg, as, ps.Logger))
		if result.Project == nil {
			result.Project = prj
		}

		if isBuild {
			if phase == "deps" {
				return runDeps(prj, result, rc, as, ps)
			}
			page := &vbuild.BuildResult{Project: prj, BuildResult: result}
			return render(rc, as, page, ps, "projects", actT.Title)
		}

		ps.Title = fmt.Sprintf("[%s] %s", actT.Title, prj.Title())
		ps.Data = result
		page := &vaction.Result{Ctx: &action.ResultContext{Prj: prj, Cfg: cfg, Res: result}, IsBuild: actT.Key == action.TypeBuild.Key}
		return render(rc, as, page, ps, "projects", prj.Key, actT.Title)
	})
}

func RunAllActions(rc *fasthttp.RequestCtx) {
	act("run.all.actions", rc, func(as *app.State, ps *cutil.PageState) (string, error) {
		actS, err := RCRequiredString(rc, "act", false)
		if err != nil {
			return "", err
		}
		cfg := util.ValueMap{}
		rc.QueryArgs().VisitAll(func(k []byte, v []byte) {
			cfg[string(k)] = string(v)
		})
		actT := action.TypeFromString(actS)
		prjs := as.Services.Projects.Projects()

		if actT.Key == action.TypeBuild.Key {
			switch cfg["phase"] {
			case nil:
				page := &vaction.Results{T: actT, Cfg: cfg, Ctxs: nil, IsBuild: true}
				return render(rc, as, page, ps, "projects", actT.Title)
			case "deps":
				return runAllDeps(cfg, prjs, rc, as, ps)
			}
		}

		results, _ := util.AsyncCollect(prjs, func(prj *project.Project) (*action.ResultContext, error) {
			c := cfg.Clone()
			result := action.Apply(ps.Context, actionParams(prj.Key, actT, c, as, ps.Logger))
			if result.Project == nil {
				result.Project = prj
			}
			return &action.ResultContext{Prj: prj, Cfg: c, Res: result}, nil
		})
		slices.SortFunc(results, func(l *action.ResultContext, r *action.ResultContext) bool {
			return strings.ToLower(l.Prj.Title()) < strings.ToLower(r.Prj.Title())
		})
		ps.Title = fmt.Sprintf("[%s] All Projects", actT.Title)
		ps.Data = results
		page := &vaction.Results{T: actT, Cfg: cfg, Ctxs: results}
		return render(rc, as, page, ps, "projects", actT.Title)
	})
}
