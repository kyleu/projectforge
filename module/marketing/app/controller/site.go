package controller

import (
	"github.com/fasthttp/router"
	"github.com/valyala/fasthttp"

	"{{{ .Package }}}/app"
	"{{{ .Package }}}/app/controller/cutil"
	"{{{ .Package }}}/app/site"
	"{{{ .Package }}}/app/telemetry/httpmetrics"
	"{{{ .Package }}}/app/util"
)

func SiteRoutes() fasthttp.RequestHandler {
	r := router.New()

	r.GET("/", Site)

	r.GET(defaultProfilePath, ProfileSite)
	r.POST(defaultProfilePath, ProfileSave)
	r.GET("/auth/{key}", AuthDetail)
	r.GET("/auth/callback/{key}", AuthCallback)
	r.GET("/auth/logout/{key}", AuthLogout)

	r.GET("/favicon.ico", Favicon)
	r.GET("/assets/{_:*}", Static)

	r.GET("/{path:*}", Site)

	r.OPTIONS("/", Options)
	r.OPTIONS("/{_:*}", Options)
	r.NotFound = NotFound

	p := httpmetrics.NewMetrics("marketing_site")
	return fasthttp.CompressHandlerBrotliLevel(p.WrapHandler(r), fasthttp.CompressBrotliBestSpeed, fasthttp.CompressBestSpeed)
}

func Site(rc *fasthttp.RequestCtx) {
	actSite("site", rc, func(as *app.State, ps *cutil.PageState) (string, error) {
		path := util.SplitAndTrim(string(rc.Request.URI().Path()), "/")
		redir, page, bc, err := site.Handle(path, rc, as, ps)
		if err != nil {
			return "", err
		}
		if redir != "" {
			return redir, nil
		}
		return render(rc, as, page, ps, bc...)
	})
}
