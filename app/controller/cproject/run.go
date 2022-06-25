package cproject

import (
	"fmt"

	"github.com/valyala/fasthttp"

	"projectforge.dev/projectforge/app"
	"projectforge.dev/projectforge/app/action"
	"projectforge.dev/projectforge/app/controller"
	"projectforge.dev/projectforge/app/controller/cutil"
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
