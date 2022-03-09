// Content managed by Project Forge, see [projectforge.md] for details.
package controller

import (
	"github.com/valyala/fasthttp"

	"projectforge.dev/projectforge/app"
	"projectforge.dev/projectforge/app/controller/cutil"
	"projectforge.dev/projectforge/app/util"
	"projectforge.dev/projectforge/views"
)

func About(rc *fasthttp.RequestCtx) {
	act("about", rc, func(as *app.State, ps *cutil.PageState) (string, error) {
		ps.Data = util.AppName + " v" + as.BuildInfo.Version
		return render(rc, as, &views.About{}, ps, "about")
	})
}
