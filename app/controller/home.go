// Package controller: $PF_IGNORE$
package controller

import (
	"github.com/kyleu/projectforge/app"
	"github.com/kyleu/projectforge/app/controller/cutil"
	"github.com/kyleu/projectforge/views"
	"github.com/pkg/errors"
	"github.com/valyala/fasthttp"
)

func Home(ctx *fasthttp.RequestCtx) {
	act("home", ctx, func(as *app.State, ps *cutil.PageState) (string, error) {
		prj, err := as.Services.Projects.Root()
		if err != nil {
			return "", errors.Wrap(err, "unable to load root project")
		}

		ps.Data = prj
		return render(ctx, as, &views.Home{}, ps)
	})
}
