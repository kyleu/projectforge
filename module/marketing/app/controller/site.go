package controller

import (
	"net/http"
	"strings"

	"{{{ .Package }}}/app"
	"{{{ .Package }}}/app/controller/cutil"
	"{{{ .Package }}}/app/site"
	"{{{ .Package }}}/app/util"
	"{{{ .Package }}}/views/verror"
)

func Site(w http.ResponseWriter, r *http.Request) {
	path := util.StringSplitAndTrim(string(r.URL.Path), "/")
	action := "site"
	if len(path) > 0 {
		action += "." + strings.Join(path, ".")
	}
	ActSite(action, w, r, func(as *app.State, ps *cutil.PageState) (string, error) {
		redir, page, bc, err := site.Handle(path, as, ps)
		if err != nil {
			return "", err
		}
		if _, ok := page.(*verror.NotFound); ok {
			w.WriteHeader(http.StatusNotFound)
		}
		if redir != "" {
			return redir, nil
		}
		return Render(w, r, as, page, ps, bc...)
	})
}
