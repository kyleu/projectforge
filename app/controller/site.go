// Content managed by Project Forge, see [projectforge.md] for details.
package controller

import (
	"github.com/fasthttp/router"
	"github.com/valyala/fasthttp"

	"projectforge.dev/app"
	"projectforge.dev/app/controller/cutil"
	"projectforge.dev/app/lib/telemetry/httpmetrics"
	"projectforge.dev/app/site"
	"projectforge.dev/app/util"
)

func SiteRoutes() fasthttp.RequestHandler {
	r := router.New()

	r.GET("/", Site)

	r.GET(defaultProfilePath, ProfileSite)
	r.POST(defaultProfilePath, ProfileSave)

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
		path := util.StringSplitAndTrim(string(rc.Request.URI().Path()), "/")
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
