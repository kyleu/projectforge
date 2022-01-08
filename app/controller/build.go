package controller

import (
	"github.com/kyleu/projectforge/app/action"
	"github.com/valyala/fasthttp"

	"github.com/kyleu/projectforge/app"
	"github.com/kyleu/projectforge/app/controller/cutil"
	"github.com/kyleu/projectforge/views/vbuild"
)

func BuildIndex(rc *fasthttp.RequestCtx) {
	act("build.index", rc, func(as *app.State, ps *cutil.PageState) (string, error) {
		prj, err := getProject(rc, as)
		if err != nil {
			return "", err
		}

		gitStatus, _ := as.Services.Git.Status(prj)
		ps.Data = action.AllBuilds
		return render(rc, as, &vbuild.BuildResult{Project: prj, GitResult: gitStatus}, ps, "projects", prj.Key, "Build")
	})
}
