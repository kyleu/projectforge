package clib

import (
	"net/http"

	"{{{ .Package }}}/app"
	"{{{ .Package }}}/app/controller"
	"{{{ .Package }}}/app/controller/cutil"
	"{{{ .Package }}}/app/util"
	"{{{ .Package }}}/views"
)

func About(w http.ResponseWriter, r *http.Request) {
	controller.Act("about", w, r, func(as *app.State, ps *cutil.PageState) (string, error) {
		ps.SetTitleAndData("About "+util.AppName, util.AppName+" v"+as.BuildInfo.Version)
		return controller.Render(w, r, as, &views.About{}, ps, "about")
	})
}
