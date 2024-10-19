// Package controller - $PF_GENERATE_ONCE$
package controller

import (
	"net/http"

	"{{{ .Package }}}/app"
	"{{{ .Package }}}/app/controller/cutil"
	"{{{ .Package }}}/app/util"
	"{{{ .Package }}}/views"
)

var homeContent = util.ValueMap{
	util.AppKey: util.AppName,
	"urls": map[string]string{
		"home": "/",
	},
}

func Home(w http.ResponseWriter, r *http.Request) {
	Act("home", w, r, func(as *app.State, ps *cutil.PageState) (string, error) {
		ps.Data = homeContent
		return Render(r, as, &views.Home{}, ps)
	})
}
