package controller

import (
	"github.com/fasthttp/router"
	"github.com/valyala/fasthttp"

	"{{{ .Package }}}/app"
	"{{{ .Package }}}/app/controller/cutil"
	"{{{ .Package }}}/app/site"
	"{{{ .Package }}}/app/util"
)

func SiteRoutes() *router.Router {
	w := fasthttp.CompressHandler
	r := router.New()

	r.GET("/", w(Site))

	r.GET(defaultProfilePath, w(ProfileSite))
	r.POST(defaultProfilePath, w(ProfileSave)){{{ if .HasModule "oauth" }}}
	r.GET("/auth/{key}", w(AuthDetail))
	r.GET("/auth/callback/{key}", w(AuthCallback))
	r.GET("/auth/logout/{key}", w(AuthLogout)){{{ end }}}

	r.GET("/favicon.ico", Favicon)
	r.GET("/assets/{_:*}", Static)

	r.GET("/{path:*}", w(Site))

	r.OPTIONS("/", w(Options))
	r.OPTIONS("/{_:*}", w(Options))
	r.NotFound = NotFound

	return r
}

func Site(ctx *fasthttp.RequestCtx) {
	actSite("site", ctx, func(as *app.State, ps *cutil.PageState) (string, error) {
		path := util.SplitAndTrim(string(ctx.Request.URI().Path()), "/")
		redir, page, bc, err := site.Handle(path, ctx, as, ps)
		if err != nil {
			return "", err
		}
		if redir != "" {
			return redir, nil
		}
		return render(ctx, as, page, ps, bc...)
	})
}
