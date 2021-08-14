package controller

import (
	"github.com/fasthttp/router"
	"github.com/valyala/fasthttp"

	"{{{ .Package }}}/app"
	"{{{ .Package }}}/app/controller/cutil"
	"{{{ .Package }}}/app/site"
	"{{{ .Package }}}/app/telemetry"
	"{{{ .Package }}}/app/util"
)

func SiteRoutes() (*telemetry.Metrics, fasthttp.RequestHandler) {
	w := fasthttp.CompressHandler
	r := router.New()

	r.GET("/", w(Site))

	r.GET(defaultProfilePath, w(ProfileSite))
	r.POST(defaultProfilePath, w(ProfileSave))
	r.GET("/auth/{key}", w(AuthDetail))
	r.GET("/auth/callback/{key}", w(AuthCallback))
	r.GET("/auth/logout/{key}", w(AuthLogout))

	r.GET("/favicon.ico", Favicon)
	r.GET("/assets/{_:*}", Static)

	r.GET("/{path:*}", w(Site))

	r.OPTIONS("/", w(Options))
	r.OPTIONS("/{_:*}", w(Options))
	r.NotFound = NotFound

	m := telemetry.NewMetrics("marketing_site")
	return m, m.WrapHandler(r)
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