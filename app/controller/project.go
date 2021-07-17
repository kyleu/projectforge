package controller

import (
	"fmt"

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

		prj, err := as.Services.Projects.Get(key)
		if err != nil {
			return "", err
		}

		ps.Title = fmt.Sprintf("%s (project %s)", prj.Name, prj.Key)
		ps.Data = prj
		return render(ctx, as, &vproject.Detail{Project: prj}, ps, "projects", prj.Key)
	})
}
