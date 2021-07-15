package controller

import (
	"runtime/debug"

	"github.com/valyala/fasthttp"

	"$PF_PACKAGE$/app"
	"$PF_PACKAGE$/app/controller/cutil"
	"$PF_PACKAGE$/views/verror"
	"$PF_PACKAGE$/views/vhelp"
	"github.com/pkg/errors"
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

func Modules(ctx *fasthttp.RequestCtx) {
	act("modules", ctx, func(as *app.State, ps *cutil.PageState) (string, error) {
		mods, ok := debug.ReadBuildInfo()
		if !ok {
			return "", errors.New("unable to gather modules")
		}
		ps.Title = "Modules"
		ps.Data = mods.Deps
		return render(ctx, as, &vhelp.Modules{Mods: mods.Deps}, ps, "modules")
	})
}
