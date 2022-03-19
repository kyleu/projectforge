package controller

import (
	"github.com/pkg/errors"
	"github.com/valyala/fasthttp"
	"projectforge.dev/projectforge/app/action"
	"projectforge.dev/projectforge/app/build"
	"projectforge.dev/projectforge/app/util"

	"projectforge.dev/projectforge/app"
	"projectforge.dev/projectforge/app/controller/cutil"
	"projectforge.dev/projectforge/views/vbuild"
)

func Build(rc *fasthttp.RequestCtx) {
	actKey, _ := RCRequiredString(rc, "act", true)
	act("build."+actKey, rc, func(as *app.State, ps *cutil.PageState) (string, error) {
		prj, err := getProject(rc, as)
		if err != nil {
			return "", err
		}

		bc := []string{"projects", prj.Key}

		var res *action.Result
		if actKey == "" {
			bc = append(bc, "Build")
			ps.Data = action.AllBuilds
		} else {
			phase := action.AllBuilds.Get(actKey)
			args := util.ValueMap{"phase": actKey}
			rc.URI().QueryArgs().VisitAll(func(key []byte, value []byte) {
				args[string(key)] = string(value)
			})
			prms := actionParams(ps.Span, prj.Key, action.TypeBuild, args, as, ps.Logger)
			res = action.Apply(ps.Context, prms)
			bc = append(bc, "Build||/b/"+prj.Key, phase.Title)
			ps.Data = res
		}

		if actKey == "deps" {
			deps, ok := res.Data.(build.Dependencies)
			if !ok {
				return "", errors.Errorf("data is of type [%T], expected [Dependencies]", res.Data)
			}
			return render(rc, as, &vbuild.Deps{Project: prj, BuildResult: res, Dependencies: deps}, ps, bc...)
		}
		gitStatus, _ := as.Services.Git.Status(prj)
		return render(rc, as, &vbuild.BuildResult{Project: prj, BuildResult: res, GitResult: gitStatus}, ps, bc...)
	})
}
