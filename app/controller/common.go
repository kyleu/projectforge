// Content managed by Project Forge, see [projectforge.md] for details.
package controller

import (
	"github.com/valyala/fasthttp"

	"projectforge.dev/app"
	"projectforge.dev/app/controller/cutil"
	"projectforge.dev/app/lib/user"
	"projectforge.dev/app/util"
	"projectforge.dev/views/verror"
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
		path := string(rc.Request.URI().Path())
		ps.Logger.Warnf("%s %s returned [%d]", string(rc.Method()), path, fasthttp.StatusNotFound)
		if ps.Title == "" {
			ps.Title = "Page not found"
		}
		bc := util.StringSplitAndTrim(string(rc.URI().Path()), "/")
		bc = append(bc, "Not Found")
		ps.Data = ps.Title
		return render(rc, as, &verror.NotFound{Path: path}, ps, bc...)
	})
}

func Unauthorized(rc *fasthttp.RequestCtx, reason string, accounts user.Accounts) func(as *app.State, ps *cutil.PageState) (string, error) {
	return func(as *app.State, ps *cutil.PageState) (string, error) {
		cutil.WriteCORS(rc)
		rc.Response.Header.Set("Content-Type", "text/html; charset=utf-8")
		rc.SetStatusCode(fasthttp.StatusUnauthorized)
		path := string(rc.Request.URI().Path())
		ps.Logger.Warnf("%s %s returned [%d]", string(rc.Method()), path, fasthttp.StatusNotFound)
		if ps.Title == "" {
			ps.Title = "Unauthorized"
		}
		bc := util.StringSplitAndTrim(string(rc.URI().Path()), "/")
		bc = append(bc, "Unauthorized")
		ps.Data = ps.Title
		return render(rc, as, &verror.Unauthorized{Path: path, Message: reason, Accounts: accounts}, ps, bc...)
	}
}
