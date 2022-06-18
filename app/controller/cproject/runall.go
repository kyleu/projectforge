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
)

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

		var results []*action.ResultContext

		results = action.ApplyAll(ps.Context, prjs, actT, cfg, as, ps.Logger)

		ps.Title = fmt.Sprintf("[%s] All Projects", actT.Title)
		ps.Data = results
		page := &vaction.Results{T: actT, Cfg: cfg, Projects: prjs, Ctxs: results, Tags: tags}
		return controller.Render(rc, as, page, ps, "projects", actT.Title)
	})
}
