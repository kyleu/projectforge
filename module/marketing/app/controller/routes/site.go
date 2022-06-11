package routes

import (
	"github.com/fasthttp/router"
	"github.com/valyala/fasthttp"

	"{{{ .Package }}}/app/controller"
	"{{{ .Package }}}/app/controller/clib"
	"{{{ .Package }}}/app/lib/telemetry/httpmetrics"
)

func SiteRoutes() fasthttp.RequestHandler {
	r := router.New()

	r.GET("/", controller.Site)

	r.GET(controller.DefaultProfilePath, clib.ProfileSite)
	r.POST(controller.DefaultProfilePath, clib.ProfileSave){{{ if .HasModule "oauth" }}}
	r.GET("/auth/{key}", clib.AuthDetail)
	r.GET("/auth/callback/{key}", clib.AuthCallback)
	r.GET("/auth/logout/{key}", clib.AuthLogout){{{ end }}}

	r.GET("/favicon.ico", clib.Favicon)
	r.GET("/assets/{_:*}", clib.Static)

	r.GET("/{path:*}", controller.Site)

	r.OPTIONS("/", controller.Options)
	r.OPTIONS("/{_:*}", controller.Options)
	r.NotFound = controller.NotFound

	p := httpmetrics.NewMetrics("marketing_site")
	return fasthttp.CompressHandlerLevel(p.WrapHandler(r), fasthttp.CompressBestSpeed)
}
