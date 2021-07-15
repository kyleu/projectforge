package controller

import (
	"$PF_PACKAGE$/app"
	"$PF_PACKAGE$/app/controller/cutil"
	"$PF_PACKAGE$/app/util"
	"$PF_PACKAGE$/views"
	"github.com/valyala/fasthttp"
)

var homeContent = util.ValueMap{
	"_": util.AppName,
	"urls": map[string]string{
		"TODO": "/todo",
	},
}

func Home(ctx *fasthttp.RequestCtx) {
	act("home", ctx, func(as *app.State, ps *cutil.PageState) (string, error) {
		ps.Data = homeContent
		return render(ctx, as, &views.Home{}, ps)
	})
}
