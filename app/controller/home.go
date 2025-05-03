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
		ps.SetTitleAndData(util.AppName, util.ValueMap{"projects": prjs, "modules": mods})
		return Render(r, as, &views.Home{Projects: prjs, Execs: execs, Modules: mods}, ps)
	})
}
