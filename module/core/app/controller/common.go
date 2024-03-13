package controller

import (
	"github.com/valyala/fasthttp"

	"{{{ .Package }}}/app"
	"{{{ .Package }}}/app/controller/cutil"{{{ if .HasAccount }}}
	"{{{ .Package }}}/app/lib/user"{{{ end }}}
	"{{{ .Package }}}/app/util"
	"{{{ .Package }}}/views/verror"
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

func Unauthorized(rc *fasthttp.RequestCtx, reason string{{{ if .HasAccount }}}, accounts user.Accounts{{{ end }}}) func(as *app.State, ps *cutil.PageState) (string, error) {
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
			{{{ if .HasAccount }}}if len(accounts) == 0 {
				reason = "not signed in"
			} else {
				reason = "no access"
			}{{{ else }}}reason = "no access"{{{ end }}}
		}
		ps.Data = util.ValueMap{"status": "unauthorized", "statusCode": fasthttp.StatusUnauthorized, "message": reason}
		return Render(rc, as, &verror.Unauthorized{Path: path, Message: reason{{{ if .HasAccount }}}, Accounts: accounts{{{ end }}}}, ps, bc...)
	}
}
