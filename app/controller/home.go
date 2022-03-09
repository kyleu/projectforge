// Package controller $PF_IGNORE$
package controller

import (
	"github.com/valyala/fasthttp"
	"projectforge.dev/projectforge/app"
	"projectforge.dev/projectforge/app/controller/cutil"
	"projectforge.dev/projectforge/views"
)

func Home(rc *fasthttp.RequestCtx) {
	act("home", rc, func(as *app.State, ps *cutil.PageState) (string, error) {
		prjs := as.Services.Projects.Projects()
		mods := as.Services.Modules.Modules()
		ps.Data = map[string]interface{}{"projects": prjs, "modules": mods}
		return render(rc, as, &views.Home{Projects: prjs, Modules: mods}, ps)
	})
}
