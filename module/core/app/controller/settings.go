package controller

import (
	"github.com/valyala/fasthttp"

	"{{{ .Package }}}/app"
	"{{{ .Package }}}/app/controller/cutil"
	"{{{ .Package }}}/app/settings"
	"{{{ .Package }}}/views/vsettings"
)

func Settings(ctx *fasthttp.RequestCtx) {
	act("settings", ctx, func(as *app.State, ps *cutil.PageState) (string, error) {
		current := &settings.Settings{}
		ps.Title = "Settings"
		ps.Data = current
		return render(ctx, as, &vsettings.Settings{Settings: current}, ps, "settings")
	})
}
