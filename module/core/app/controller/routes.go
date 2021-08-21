package controller

import (
	"github.com/fasthttp/router"
	"github.com/valyala/fasthttp"

	"{{{ .Package }}}/app/telemetry/httpmetrics"
	"{{{ .Package }}}/app/util"
)

//nolint
func AppRoutes() fasthttp.RequestHandler {
	r := router.New()

	r.GET("/", Home)
	r.GET("/about", About)
	r.GET("/settings", Settings)
	r.GET("/theme", ThemeList)
	r.GET("/theme/{key}", ThemeEdit)
	r.POST("/theme/{key}", ThemeSave)
	r.GET(defaultSearchPath, Search)

	r.GET(defaultProfilePath, Profile)
	r.POST(defaultProfilePath, ProfileSave)
	r.GET("/auth/{key}", AuthDetail)
	r.GET("/auth/callback/{key}", AuthCallback)
	r.GET("/auth/logout/{key}", AuthLogout)

	r.GET("/admin", Admin)
	r.GET("/admin/{path:*}", Admin)

	// $PF_SECTION_START(routes)$
	// $PF_SECTION_END(routes)$

	r.GET("/sandbox", SandboxList)
	r.GET("/sandbox/{key}", SandboxRun)

	r.GET("/favicon.ico", Favicon)
	r.GET("/robots.txt", RobotsTxt)
	r.GET("/assets/{_:*}", Static)

	r.OPTIONS("/", Options)
	r.OPTIONS("/{_:*}", Options)
	r.NotFound = NotFound

	p := httpmetrics.NewMetrics(util.AppKey)
	return fasthttp.CompressHandler(p.WrapHandler(r))
}
