package controller

import (
	"github.com/kyleu/projectforge/views"
	"github.com/kyleu/projectforge/views/vproject"
	"github.com/pkg/errors"
	"github.com/valyala/fasthttp"

	"github.com/kyleu/projectforge/app/controller/cutil"

	"github.com/kyleu/projectforge/app"
)

func RootProjectDetail(ctx *fasthttp.RequestCtx) {
	act("project.root", ctx, func(as *app.State, ps *cutil.PageState) (string, error) {
		prj, err := as.Services.Projects.Root()
		if err != nil {
			return "", errors.Wrap(err, "unable to load root project")
		}
		results := prj
		ps.Title = "Project [" + prj.Key + "]"
		ps.Data = results
		return render(ctx, as, &vproject.Detail{Project: prj}, ps, "root")
	})
}

func ProjectDetail(ctx *fasthttp.RequestCtx) {
	act("project.detail", ctx, func(as *app.State, ps *cutil.PageState) (string, error) {
		key, err := ctxRequiredString(ctx, "key", true)
		if err != nil {
			return "", err
		}

		results := key

		ps.Title = "Project [" + key + "]"
		ps.Data = results
		return render(ctx, as, &views.Debug{}, ps, "project", key)
	})
}
