// Package controller: $PF_IGNORE$
package controller

import (
	"github.com/fasthttp/router"
	"github.com/valyala/fasthttp"
)

//nolint
func AppRoutes() *router.Router {
	w := fasthttp.CompressHandler
	r := router.New()

	r.GET("/", w(Home))
	r.GET("/settings", w(Settings))
	r.GET("/theme", w(ThemeList))
	r.GET("/theme/{key}", w(ThemeEdit))
	r.POST("/theme/{key}", w(ThemeSave))
	r.GET(defaultSearchPath, w(Search))

	r.GET(defaultProfilePath, w(Profile))
	r.POST(defaultProfilePath, w(ProfileSave))
	r.GET("/auth/{key}", w(AuthDetail))
	r.GET("/auth/{key}/callback", w(AuthCallback))
	r.GET("/auth/{key}/logout", w(AuthLogout))

	r.GET("/help", w(Help))
	r.GET("/help/{path:*}", w(Help))

	// $PF_SECTION(routes)
	r.GET("/run", w(Run))
	r.GET("/run/{tgt}/{act}", w(RunAction))

	r.GET("/tools", w(ToolList))
	r.GET("/tools/svg", w(SVGList))
	r.GET("/tools/svg/add", w(SVGAdd))
	r.GET("/tools/svg/build", w(SVGBuild))
	// $PF_SECTION(routes)

	r.GET("/sandbox", w(SandboxList))
	r.GET("/sandbox/{key}", w(SandboxRun))

	r.GET("/modules", w(Modules))

	r.GET("/favicon.ico", Favicon)
	r.GET("/robots.txt", RobotsTxt)
	r.GET("/assets/{_:*}", Static)

	r.OPTIONS("/", w(Options))
	r.OPTIONS("/{_:*}", w(Options))
	r.NotFound = NotFound

	return r
}
