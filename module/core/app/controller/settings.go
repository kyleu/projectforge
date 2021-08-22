package controller

import (
	"github.com/valyala/fasthttp"

	"{{{ .Package }}}/app"
	"{{{ .Package }}}/app/controller/cutil"
	"{{{ .Package }}}/app/settings"
	"{{{ .Package }}}/views/vsettings"
)

func Settings(rc *fasthttp.RequestCtx) {
	act("settings", rc, func(as *app.State, ps *cutil.PageState) (string, error) {
		current := &settings.Settings{}
		ps.Title = "Settings"
		ps.Data = current
		return render(rc, as, &vsettings.Settings{Settings: current}, ps, "settings")
	})
}
