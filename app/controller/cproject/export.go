package cproject

import (
	"fmt"
	"net/http"

	"projectforge.dev/projectforge/app"
	"projectforge.dev/projectforge/app/controller"
	"projectforge.dev/projectforge/app/controller/cutil"
	"projectforge.dev/projectforge/views/vexport"
)

func ProjectExportOverview(w http.ResponseWriter, r *http.Request) {
	controller.Act("project.export.overview", w, r, func(as *app.State, ps *cutil.PageState) (string, error) {
		prj, err := getProject(r, as)
		if err != nil {
			return "", err
		}
		args, err := prj.ModuleArgExport(as.Services.Projects, ps.Logger)
		if err != nil {
			return "", err
		}
		bc := []string{"projects", prj.Key, "Export"}
		ps.SetTitleAndData(fmt.Sprintf("[%s] Export", prj.Key), args)
		return controller.Render(w, r, as, &vexport.Overview{Project: prj, Args: args}, ps, bc...)
	})
}
