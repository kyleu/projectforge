package controller

import (
	"github.com/fasthttp/router"
	"github.com/valyala/fasthttp"

	"{{{ .Package }}}/app/lib/telemetry/httpmetrics"
	"{{{ .Package }}}/app/util"
)

//nolint
func AppRoutes() fasthttp.RequestHandler {
	r := router.New()

	r.GET("/", Home)
	r.GET("/healthcheck", Healthcheck)
	r.GET("/about", About)
	r.GET("/theme", ThemeList)
	r.GET("/theme/{key}", ThemeEdit)
	r.POST("/theme/{key}", ThemeSave){{{ if .HasModule "search" }}}
	r.GET(defaultSearchPath, Search){{{ end }}}

	r.GET(defaultProfilePath, Profile)
	r.POST(defaultProfilePath, ProfileSave){{{ if .HasModule "oauth" }}}
	r.GET("/auth/{key}", AuthDetail)
	r.GET("/auth/callback/{key}", AuthCallback)
	r.GET("/auth/logout/{key}", AuthLogout){{{ end }}}{{{ if.HasModule "oauth" }}}

	// $PF_INJECT_START(codegen)$
	// $PF_INJECT_END(codegen)${{{ end }}}

	// $PF_SECTION_START(routes)$
	// $PF_SECTION_END(routes)$

	r.GET("/admin", Admin)
	r.GET("/admin/sandbox", SandboxList)
	r.GET("/admin/sandbox/{key}", SandboxRun)
	r.GET("/admin/{path:*}", Admin)

	r.GET("/favicon.ico", Favicon)
	r.GET("/robots.txt", RobotsTxt)
	r.GET("/assets/{_:*}", Static)

	r.OPTIONS("/", Options)
	r.OPTIONS("/{_:*}", Options)
	r.NotFound = NotFound

	p := httpmetrics.NewMetrics(util.AppKey)
	return fasthttp.CompressHandlerBrotliLevel(p.WrapHandler(r), fasthttp.CompressBrotliBestSpeed, fasthttp.CompressBestSpeed)
}
