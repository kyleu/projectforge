package controller

import (
	"net/http"

	"{{{ .Package }}}/app"
	"{{{ .Package }}}/app/controller/cutil"
	"{{{ .Package }}}/app/site"
	"{{{ .Package }}}/app/util"
	"{{{ .Package }}}/views/verror"
)

func Site(w http.ResponseWriter, r *http.Request) {
	path := util.StringSplitAndTrim(r.URL.Path, "/")
	action := "site"
	if len(path) > 0 {
		action += "." + util.StringJoin(path, ".")
	}
	ActSite(action, w, r, func(as *app.State, ps *cutil.PageState) (string, error) {
		redir, page, bc, err := site.Handle(as, path, ps)
		if err != nil {
			return "", err
		}
		if _, err := util.Cast[*verror.NotFound](page); err != nil {
			w.WriteHeader(http.StatusNotFound)
		}
		if redir != "" {
			return redir, nil
		}
		return Render(r, as, page, ps, bc...)
	})
}
