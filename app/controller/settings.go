package controller

import (
	"github.com/valyala/fasthttp"

	"github.com/kyleu/projectforge/app/controller/cutil"
	"github.com/kyleu/projectforge/app/settings"

	"github.com/kyleu/projectforge/app"
	"github.com/kyleu/projectforge/views/vsettings"
)

func Settings(ctx *fasthttp.RequestCtx) {
	act("settings", ctx, func(as *app.State, ps *cutil.PageState) (string, error) {
		current := &settings.Settings{}
		ps.Title = "Settings"
		ps.Data = current
		return render(ctx, as, &vsettings.Settings{Settings: current}, ps, "settings")
	})
}
