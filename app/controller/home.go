// Package controller $PF_IGNORE$
package controller

import (
	"github.com/valyala/fasthttp"

	"projectforge.dev/projectforge/app"
	"projectforge.dev/projectforge/app/controller/cutil"
	"projectforge.dev/projectforge/app/util"
	"projectforge.dev/projectforge/views"
)

func Home(rc *fasthttp.RequestCtx) {
	Act("home", rc, func(as *app.State, ps *cutil.PageState) (string, error) {
		prjs := as.Services.Projects.Projects()
		mods := as.Services.Modules.Modules()
		ps.Data = util.ValueMap{"projects": prjs, "modules": mods}
		return Render(rc, as, &views.Home{Projects: prjs, Modules: mods}, ps)
	})
}
