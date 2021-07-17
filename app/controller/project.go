package controller

import (
	"fmt"

	"github.com/kyleu/projectforge/views/vproject"
	"github.com/valyala/fasthttp"

	"github.com/kyleu/projectforge/app/controller/cutil"

	"github.com/kyleu/projectforge/app"
)

func ProjectList(ctx *fasthttp.RequestCtx) {
	act("project.root", ctx, func(as *app.State, ps *cutil.PageState) (string, error) {
		prjs := as.Services.Projects.Projects()
		ps.Title = "Project Listing"
		ps.Data = prjs
		return render(ctx, as, &vproject.List{Projects: prjs}, ps, "projects")
	})
}

func ProjectDetail(ctx *fasthttp.RequestCtx) {
	act("project.detail", ctx, func(as *app.State, ps *cutil.PageState) (string, error) {
		key, err := ctxRequiredString(ctx, "key", true)
		if err != nil {
			return "", err
		}

		prj, err := as.Services.Projects.Get(key)
		if err != nil {
			return "", err
		}

		ps.Title = fmt.Sprintf("%s (project %s)", prj.Title(), prj.Key)
		ps.Data = prj
		return render(ctx, as, &vproject.Detail{Project: prj}, ps, "projects", prj.Key)
	})
}
