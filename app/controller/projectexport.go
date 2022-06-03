package controller

import (
	"fmt"

	"github.com/valyala/fasthttp"

	"projectforge.dev/projectforge/app"
	"projectforge.dev/projectforge/app/controller/cutil"
	"projectforge.dev/projectforge/views/vexport"
)

func ProjectExportModelDetail(rc *fasthttp.RequestCtx) {
	act("project.export.model.detail", rc, func(as *app.State, ps *cutil.PageState) (string, error) {
		prj, err := getProject(rc, as)
		if err != nil {
			return "", err
		}

		modelKey, err := cutil.RCRequiredString(rc, "model", false)
		if err != nil {
			return "", err
		}

		args, err := prj.Info.ModuleArgExport()
		if err != nil {
			return "", err
		}

		model := args.Models.Get(modelKey)

		ps.Data = model

		bc := []string{"projects", prj.Key, "Export||/p/" + prj.Key, modelKey}
		ps.Title = fmt.Sprintf("[%s] %s", prj.Key, modelKey)
		return render(rc, as, &vexport.ModelDetail{Model: model}, ps, bc...)
	})
}
