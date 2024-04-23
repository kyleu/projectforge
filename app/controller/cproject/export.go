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
		prj, err := getProjectWithArgs(r, as, ps.Logger)
		if err != nil {
			return "", err
		}
		bc := []string{"projects", prj.Key, "Export"}
		ps.SetTitleAndData(fmt.Sprintf("[%s] Export", prj.Key), prj.ExportArgs)
		return controller.Render(r, as, &vexport.Overview{Project: prj, Args: prj.ExportArgs}, ps, bc...)
	})
}
