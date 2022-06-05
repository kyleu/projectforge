package controller

import (
	"fmt"

	"github.com/valyala/fasthttp"
	"projectforge.dev/projectforge/app"
	"projectforge.dev/projectforge/app/controller/cutil"
	"projectforge.dev/projectforge/views/vexport"
)

func ProjectExportOverview(rc *fasthttp.RequestCtx) {
	act("project.export.overview", rc, func(as *app.State, ps *cutil.PageState) (string, error) {
		prj, err := getProject(rc, as)
		if err != nil {
			return "", err
		}
		args, err := prj.Info.ModuleArgExport()
		if err != nil {
			return "", err
		}
		ps.Data = args

		bc := []string{"projects", prj.Key, "Export"}
		ps.Title = fmt.Sprintf("[%s] Export", prj.Key)
		return render(rc, as, &vexport.Overview{Project: prj, Args: args}, ps, bc...)
	})
}
