package cproject

import (
	"fmt"
	"strings"

	"github.com/valyala/fasthttp"
	"golang.org/x/exp/slices"

	"projectforge.dev/projectforge/app"
	"projectforge.dev/projectforge/app/action"
	"projectforge.dev/projectforge/app/controller"
	"projectforge.dev/projectforge/app/controller/cutil"
	"projectforge.dev/projectforge/app/project"
	"projectforge.dev/projectforge/app/util"
	"projectforge.dev/projectforge/views/vaction"
	"projectforge.dev/projectforge/views/vbuild"
)

const (
	depsKey = "deps"
	pkgsKey = "packages"
)

func RunAction(rc *fasthttp.RequestCtx) {
	controller.Act("run.action", rc, func(as *app.State, ps *cutil.PageState) (string, error) {
		tgt, err := cutil.RCRequiredString(rc, "key", false)
		if err != nil {
			return "", err
		}
		actS, err := cutil.RCRequiredString(rc, "act", false)
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
			page := &vbuild.BuildResult{Project: prj, Cfg: cfg, BuildResult: nil, GitResult: nil}
			return controller.Render(rc, as, page, ps, "projects", prj.Key, actT.Title)
		}

		result := action.Apply(ps.Context, actionParams(tgt, actT, cfg, as, ps.Logger))
		if result == nil {
			result = &action.Result{}
		}
		if result.Project == nil {
			result.Project = prj
		}

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

		ps.Title = fmt.Sprintf("[%s] %s", actT.Title, prj.Title())
		ps.Data = result
		page := &vaction.Result{Ctx: &action.ResultContext{Prj: prj, Cfg: cfg, Res: result}, IsBuild: actT.Key == action.TypeBuild.Key}
		return controller.Render(rc, as, page, ps, "projects", prj.Key, actT.Title)
	})
}

func RunAllActions(rc *fasthttp.RequestCtx) {
	controller.Act("run.all.actions", rc, func(as *app.State, ps *cutil.PageState) (string, error) {
		actS, err := cutil.RCRequiredString(rc, "act", false)
		if err != nil {
			return "", err
		}
		cfg := util.ValueMap{}
		rc.QueryArgs().VisitAll(func(k []byte, v []byte) {
			cfg[string(k)] = string(v)
		})
		actT := action.TypeFromString(actS)
		prjs := as.Services.Projects.Projects()
		tags := util.StringSplitAndTrim(string(rc.URI().QueryArgs().Peek("tags")), ",")
		if len(tags) == 0 {
			prjs = prjs.WithoutTags("all-skip")
		} else {
			prjs = prjs.WithTags(tags...)
		}

		if actT.Key == action.TypeBuild.Key {
			switch cfg["phase"] {
			case nil:
				ps.Title = "Build All Projects"
				page := &vaction.Results{T: actT, Cfg: cfg, Projects: prjs, Ctxs: nil, Tags: tags, IsBuild: true}
				return controller.Render(rc, as, page, ps, "projects", actT.Title)
			case depsKey:
				return runAllDeps(cfg, prjs, rc, as, ps)
			case pkgsKey:
				return runAllPkgs(cfg, prjs, rc, as, ps)
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
		page := &vaction.Results{T: actT, Cfg: cfg, Projects: prjs, Ctxs: results, Tags: tags}
		return controller.Render(rc, as, page, ps, "projects", actT.Title)
	})
}
