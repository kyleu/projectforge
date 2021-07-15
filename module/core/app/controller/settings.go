package controller

import (
	"github.com/valyala/fasthttp"

	"$PF_PACKAGE$/app/controller/cutil"
	"$PF_PACKAGE$/app/settings"

	"$PF_PACKAGE$/app"
	"$PF_PACKAGE$/views/vsettings"
)

func Settings(ctx *fasthttp.RequestCtx) {
	act("settings", ctx, func(as *app.State, ps *cutil.PageState) (string, error) {
		current := &settings.Settings{}
		ps.Title = "Settings"
		ps.Data = current
		return render(ctx, as, &vsettings.Settings{Settings: current}, ps, "settings")
	})
}
