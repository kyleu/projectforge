// Content managed by Project Forge, see [projectforge.md] for details.
package controller

import (
	"github.com/valyala/fasthttp"
	"projectforge.dev/projectforge/app/lib/upgrade"

	"projectforge.dev/projectforge/app"
	"projectforge.dev/projectforge/app/controller/cutil"
	"projectforge.dev/projectforge/app/util"
	"projectforge.dev/projectforge/views"
)

func About(rc *fasthttp.RequestCtx) {
	act("about", rc, func(as *app.State, ps *cutil.PageState) (string, error) {
		ps.Title = "About " + util.AppName
		ps.Data = util.AppName + " v" + as.BuildInfo.Version

		upsvc := upgrade.NewService(ps.Logger)
		err := upsvc.UpgradeIfNeeded(ps.Context, as.BuildInfo.Version, as.BuildInfo.Version, false)
		if err != nil {
			return "", err
		}

		return render(rc, as, &views.About{}, ps, "about")
	})
}
