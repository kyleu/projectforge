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
	r.GET("/about", w(About))
	r.GET("/settings", w(Settings))
	r.GET("/theme", w(ThemeList))
	r.GET("/theme/{key}", w(ThemeEdit))
	r.POST("/theme/{key}", w(ThemeSave))
	r.GET(defaultSearchPath, w(Search))

	r.GET(defaultProfilePath, w(Profile))
	r.POST(defaultProfilePath, w(ProfileSave)){{{ if .HasModule "oauth" }}}
	r.GET("/auth/{key}", w(AuthDetail))
	r.GET("/auth/callback/{key}", w(AuthCallback))
	r.GET("/auth/logout/{key}", w(AuthLogout)){{{ end }}}

	r.GET("/admin", w(Admin))
	r.GET("/admin/{path:*}", w(Admin))

	// $PF_SECTION_START(routes)$
	// $PF_SECTION_END(routes)$

	r.GET("/sandbox", w(SandboxList))
	r.GET("/sandbox/{key}", w(SandboxRun))

	r.GET("/favicon.ico", Favicon)
	r.GET("/robots.txt", RobotsTxt)
	r.GET("/assets/{_:*}", Static)

	r.OPTIONS("/", w(Options))
	r.OPTIONS("/{_:*}", w(Options))
	r.NotFound = NotFound

	return r
}
