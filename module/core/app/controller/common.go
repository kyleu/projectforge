package controller

import (
	"github.com/valyala/fasthttp"

	"{{{ .Package }}}/app"
	"{{{ .Package }}}/app/controller/cutil"
	"{{{ .Package }}}/views/verror"
)

func Options(ctx *fasthttp.RequestCtx) {
	cutil.WriteCORS(ctx)
	ctx.SetStatusCode(fasthttp.StatusOK)
}

func NotFound(ctx *fasthttp.RequestCtx) {
	act("notfound", ctx, func(as *app.State, ps *cutil.PageState) (string, error) {
		cutil.WriteCORS(ctx)
		ctx.Response.Header.Set("Content-Type", "text/html; charset=utf-8")
		ctx.SetStatusCode(fasthttp.StatusNotFound)
		ps.Logger.Warnf("%s %s returned [%d]", string(ctx.Method()), string(ctx.Request.URI().Path()), fasthttp.StatusNotFound)
		if ps.Title == "" {
			ps.Title = "404"
		}
		ps.Data = "404 not found"
		return render(ctx, as, &verror.NotFound{}, ps, "Not Found")
	})
}
