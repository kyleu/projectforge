package controller

import (
	"github.com/valyala/fasthttp"

	"$PF_PACKAGE$/app"
	"$PF_PACKAGE$/app/controller/cutil"
	"$PF_PACKAGE$/app/util"
	"$PF_PACKAGE$/views"
)

func About(ctx *fasthttp.RequestCtx) {
	act("about", ctx, func(as *app.State, ps *cutil.PageState) (string, error) {
		ps.Data = util.AppName + " v" + as.BuildInfo.Version
		return render(ctx, as, &views.About{}, ps)
	})
}
