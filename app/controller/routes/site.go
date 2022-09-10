// Content managed by Project Forge, see [projectforge.md] for details.
package routes

import (
	"github.com/fasthttp/router"
	"github.com/valyala/fasthttp"

	"projectforge.dev/projectforge/app/controller"
	"projectforge.dev/projectforge/app/controller/clib"
	"projectforge.dev/projectforge/app/controller/cutil"
	"projectforge.dev/projectforge/app/lib/telemetry/httpmetrics"
	"projectforge.dev/projectforge/app/util"
)

func SiteRoutes(logger util.Logger) fasthttp.RequestHandler {
	r := router.New()

	r.GET("/", controller.Site)

	r.GET(cutil.DefaultProfilePath, clib.ProfileSite)
	r.POST(cutil.DefaultProfilePath, clib.ProfileSave)

	r.GET("/favicon.ico", clib.Favicon)
	r.GET("/assets/{_:*}", clib.Static)

	r.GET("/{path:*}", controller.Site)

	r.OPTIONS("/", controller.Options)
	r.OPTIONS("/{_:*}", controller.Options)
	r.NotFound = controller.NotFound

	p := httpmetrics.NewMetrics("marketing_site", logger)
	return fasthttp.CompressHandlerLevel(p.WrapHandler(r), fasthttp.CompressBestSpeed)
}
