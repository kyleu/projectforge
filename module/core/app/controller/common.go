package controller

import (
	"github.com/valyala/fasthttp"

	"{{{ .Package }}}/app"
	"{{{ .Package }}}/app/controller/cutil"
	"{{{ .Package }}}/views/verror"
)

func Options(rc *fasthttp.RequestCtx) {
	cutil.WriteCORS(rc)
	rc.SetStatusCode(fasthttp.StatusOK)
}

func NotFound(rc *fasthttp.RequestCtx) {
	act("notfound", rc, func(as *app.State, ps *cutil.PageState) (string, error) {
		cutil.WriteCORS(rc)
		rc.Response.Header.Set("Content-Type", "text/html; charset=utf-8")
		rc.SetStatusCode(fasthttp.StatusNotFound)
		ps.Logger.Warnf("%s %s returned [%d]", string(rc.Method()), string(rc.Request.URI().Path()), fasthttp.StatusNotFound)
		if ps.Title == "" {
			ps.Title = "404"
		}
		ps.Data = "404 not found"
		return render(rc, as, &verror.NotFound{}, ps, "Not Found")
	})
}
