// Package controller - Content managed by Project Forge, see [projectforge.md] for details.
package controller

import (
	"github.com/valyala/fasthttp"

	"projectforge.dev/projectforge/app"
	"projectforge.dev/projectforge/app/controller/cutil"
	"projectforge.dev/projectforge/app/util"
	"projectforge.dev/projectforge/views/verror"
)

func Options(rc *fasthttp.RequestCtx) {
	cutil.WriteCORS(rc)
	rc.SetStatusCode(fasthttp.StatusOK)
}

func NotFoundAction(rc *fasthttp.RequestCtx) {
	Act("notfound", rc, func(as *app.State, ps *cutil.PageState) (string, error) {
		return NotFoundResponse(rc)(as, ps)
	})
}

func NotFoundResponse(rc *fasthttp.RequestCtx) func(as *app.State, ps *cutil.PageState) (string, error) {
	return func(as *app.State, ps *cutil.PageState) (string, error) {
		cutil.WriteCORS(rc)
		rc.SetStatusCode(fasthttp.StatusNotFound)
		path := string(rc.Request.URI().Path())
		ps.Logger.Warnf("%s %s returned [%d]", string(rc.Method()), path, fasthttp.StatusNotFound)
		if ps.Title == "" {
			ps.Title = "Page not found"
		}
		ps.Data = util.ValueMap{"status": "notfound", "statusCode": fasthttp.StatusNotFound, "message": ps.Title}
		bc := util.StringSplitAndTrim(string(rc.URI().Path()), "/")
		bc = append(bc, "Not Found")
		return Render(rc, as, &verror.NotFound{Path: path}, ps, bc...)
	}
}

func Unauthorized(rc *fasthttp.RequestCtx, reason string) func(as *app.State, ps *cutil.PageState) (string, error) {
	return func(as *app.State, ps *cutil.PageState) (string, error) {
		cutil.WriteCORS(rc)
		rc.SetStatusCode(fasthttp.StatusUnauthorized)
		path := string(rc.Request.URI().Path())
		ps.Logger.Warnf("%s %s returned [%d]", string(rc.Method()), path, fasthttp.StatusNotFound)
		bc := util.StringSplitAndTrim(string(rc.URI().Path()), "/")
		bc = append(bc, "Unauthorized")
		if ps.Title == "" {
			ps.Title = "Unauthorized"
		}
		if reason == "" {
			reason = "no access"
		}
		ps.Data = util.ValueMap{"status": "unauthorized", "statusCode": fasthttp.StatusUnauthorized, "message": reason}
		return Render(rc, as, &verror.Unauthorized{Path: path, Message: reason}, ps, bc...)
	}
}
