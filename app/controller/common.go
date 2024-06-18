package controller

import (
	"net/http"

	"projectforge.dev/projectforge/app"
	"projectforge.dev/projectforge/app/controller/cutil"
	"projectforge.dev/projectforge/app/util"
	"projectforge.dev/projectforge/views/verror"
)

func Options(w http.ResponseWriter, _ *http.Request) {
	cutil.WriteCORS(w)
	w.WriteHeader(http.StatusOK)
}

func NotFoundAction(w http.ResponseWriter, r *http.Request) {
	Act("notfound", w, r, func(as *app.State, ps *cutil.PageState) (string, error) {
		return NotFoundResponse(ps.W, r)(as, ps)
	})
}

func NotFoundResponse(w http.ResponseWriter, r *http.Request) func(as *app.State, ps *cutil.PageState) (string, error) {
	return func(as *app.State, ps *cutil.PageState) (string, error) {
		cutil.WriteCORS(w)
		w.WriteHeader(http.StatusNotFound)
		ps.Logger.Warnf("%s %s returned [%d]", r.Method, r.URL.Path, http.StatusNotFound)
		if ps.Title == "" {
			ps.Title = "Page not found"
		}
		ps.Data = util.ValueMap{"status": "notfound", "statusCode": http.StatusNotFound, "message": ps.Title}
		bc := util.StringSplitAndTrim(r.URL.Path, "/")
		bc = append(bc, "Not Found")
		return Render(r, as, &verror.NotFound{Path: r.URL.Path}, ps, bc...)
	}
}

func Unauthorized(w http.ResponseWriter, r *http.Request, reason string) func(as *app.State, ps *cutil.PageState) (string, error) {
	return func(as *app.State, ps *cutil.PageState) (string, error) {
		cutil.WriteCORS(w)
		w.WriteHeader(http.StatusUnauthorized)
		path := r.URL.Path
		ps.Logger.Warnf("%s %s returned [%d]", r.Method, path, http.StatusUnauthorized)
		bc := util.StringSplitAndTrim(r.URL.Path, "/")
		bc = append(bc, "Unauthorized")
		if ps.Title == "" {
			ps.Title = "Unauthorized"
		}
		if reason == "" {
			reason = "no access"
		}
		ps.Data = util.ValueMap{"status": "unauthorized", "statusCode": http.StatusUnauthorized, "message": reason}
		return Render(r, as, &verror.Unauthorized{Path: path, Message: reason}, ps, bc...)
	}
}
