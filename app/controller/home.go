package controller

import (
	"github.com/kyleu/projectforge/app"
	"github.com/kyleu/projectforge/app/controller/cutil"
	"github.com/kyleu/projectforge/views"
	"github.com/valyala/fasthttp"
)

func Home(ctx *fasthttp.RequestCtx) {
	act("home", ctx, func(as *app.State, ps *cutil.PageState) (string, error) {
		prjs := as.Services.Projects.Projects()
		ps.Data = prjs
		return render(ctx, as, &views.Home{Projects: prjs, Modules: as.Services.Modules.Modules()}, ps)
	})
}
