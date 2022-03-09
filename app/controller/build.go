package controller

import (
	"github.com/valyala/fasthttp"
	"projectforge.dev/projectforge/app/action"
	"projectforge.dev/projectforge/app/util"

	"projectforge.dev/projectforge/app"
	"projectforge.dev/projectforge/app/controller/cutil"
	"projectforge.dev/projectforge/views/vbuild"
)

func Build(rc *fasthttp.RequestCtx) {
	act("build.index", rc, func(as *app.State, ps *cutil.PageState) (string, error) {
		prj, err := getProject(rc, as)
		if err != nil {
			return "", err
		}

		var res *action.Result
		if act, _ := RCRequiredString(rc, "act", true); act == "" {
			ps.Data = action.AllBuilds
		} else {
			cfg := util.ValueMap{"phase": act}
			res = action.Apply(ps.Context, actionParams(ps.Span, prj.Key, action.TypeBuild, cfg, as, ps.Logger))
			ps.Data = res
		}

		gitStatus, _ := as.Services.Git.Status(prj)
		return render(rc, as, &vbuild.BuildResult{Project: prj, BuildResult: res, GitResult: gitStatus}, ps, "projects", prj.Key, "Build")
	})
}
