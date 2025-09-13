package cexport

import (
	"fmt"
	"net/http"

	"projectforge.dev/projectforge/app"
	"projectforge.dev/projectforge/app/controller"
	"projectforge.dev/projectforge/app/controller/cproject"
	"projectforge.dev/projectforge/app/controller/cutil"
	"projectforge.dev/projectforge/views/vexport"
)

const keyNew = "new"

func ProjectExportOverview(w http.ResponseWriter, r *http.Request) {
	controller.Act("project.export.overview", w, r, func(as *app.State, ps *cutil.PageState) (string, error) {
		prj, err := cproject.GetProjectWithArgs(r, as, ps.Logger)
		if err != nil {
			return "", err
		}
		bc := []string{"projects", prj.Key, "Export"}
		ps.SetTitleAndData(fmt.Sprintf("[%s] Export", prj.Key), prj.ExportArgs)
		return controller.Render(r, as, &vexport.Overview{Project: prj, Args: prj.ExportArgs}, ps, bc...)
	})
}
