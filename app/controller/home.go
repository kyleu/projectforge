package controller

import (
	"net/http"

	"projectforge.dev/projectforge/app"
	"projectforge.dev/projectforge/app/controller/cutil"
	"projectforge.dev/projectforge/app/util"
	"projectforge.dev/projectforge/views"
)

func Home(w http.ResponseWriter, r *http.Request) {
	Act("home", w, r, func(as *app.State, ps *cutil.PageState) (string, error) {
		prjs := as.Services.Projects.Projects()
		execs := as.Services.Exec.Execs
		mods := as.Services.Modules.ModulesVisible()
		ps.SetTitleAndData(util.AppName, util.ValueMap{"projects": prjs, "modules": mods, "execs": execs})
		page := &views.Home{Projects: prjs, Modules: mods, Execs: execs}
		return Render(r, as, page, ps)
	})
}
