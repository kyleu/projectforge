package controller

import (
	"$PF_PACKAGE$/app/util"
	"github.com/valyala/fasthttp"

	"$PF_PACKAGE$/app/controller/cutil"

	"$PF_PACKAGE$/app"
	"$PF_PACKAGE$/views/vhelp"
)

var helpContent = util.ValueMap{
	"_": "help",
	"urls": map[string]string{
		"home": "/",
	},
}

func Help(ctx *fasthttp.RequestCtx) {
	act("help", ctx, func(as *app.State, ps *cutil.PageState) (string, error) {
		ps.Title = "Help"
		ps.Data = helpContent
		return render(ctx, as, &vhelp.Help{}, ps, "help")
	})
}
