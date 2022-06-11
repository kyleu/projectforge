package cproject

import (
	"fmt"

	"github.com/valyala/fasthttp"

	"projectforge.dev/projectforge/app"
	"projectforge.dev/projectforge/app/controller"
	"projectforge.dev/projectforge/app/controller/cutil"
	"projectforge.dev/projectforge/views/vexport"
)

func ProjectExportOverview(rc *fasthttp.RequestCtx) {
	controller.Act("project.export.overview", rc, func(as *app.State, ps *cutil.PageState) (string, error) {
		prj, err := getProject(rc, as)
		if err != nil {
			return "", err
		}
		args, err := prj.ModuleArgExport(as.Services.Projects, ps.Logger)
		if err != nil {
			return "", err
		}
		ps.Data = args

		bc := []string{"projects", prj.Key, "Export"}
		ps.Title = fmt.Sprintf("[%s] Export", prj.Key)
		return controller.Render(rc, as, &vexport.Overview{Project: prj, Args: args}, ps, bc...)
	})
}
