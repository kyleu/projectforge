package controller

import (
	"github.com/kyleu/projectforge/app/util"
	"github.com/valyala/fasthttp"

	"github.com/kyleu/projectforge/app/controller/cutil"

	"github.com/kyleu/projectforge/app"
	"github.com/kyleu/projectforge/views/vhelp"
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
