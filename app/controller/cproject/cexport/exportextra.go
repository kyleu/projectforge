package cexport

import (
	"fmt"
	"net/http"

	"github.com/pkg/errors"

	"projectforge.dev/projectforge/app"
	"projectforge.dev/projectforge/app/controller"
	"projectforge.dev/projectforge/app/controller/cproject"
	"projectforge.dev/projectforge/app/controller/cutil"
	"projectforge.dev/projectforge/app/lib/metamodel/model"
	"projectforge.dev/projectforge/app/project"
	"projectforge.dev/projectforge/app/util"
	"projectforge.dev/projectforge/views/vexport"
)

func ProjectExportExtraTypeDetail(w http.ResponseWriter, r *http.Request) {
	controller.Act("project.export.extra.type.detail", w, r, func(as *app.State, ps *cutil.PageState) (string, error) {
		prj, t, err := exportLoadExtraType(r, as, ps.Logger)
		if err != nil {
			return "", err
		}
		bc := []string{"projects", prj.Key, fmt.Sprintf("Export||/p/%s/export", prj.Key), t.Title()}
		ps.SetTitleAndData(fmt.Sprintf("[%s] %s", prj.Key, t.Title()), t)
		return controller.Render(r, as, &vexport.ExtraTypeDetail{BaseURL: prj.WebPathModels(), T: t}, ps, bc...)
	})
}

func exportLoadExtraType(r *http.Request, as *app.State, logger util.Logger) (*project.Project, *model.Model, error) {
	prj, err := cproject.GetProjectWithArgs(r, as, logger)
	if err != nil {
		return nil, nil, err
	}
	modelKey, err := cutil.PathString(r, "t", false)
	if err != nil {
		return nil, nil, err
	}
	t := prj.ExportArgs.ExtraTypes.Get(modelKey)
	if t == nil {
		return nil, nil, errors.Errorf("no extra type found with key [%s]", modelKey)
	}
	return prj, t, nil
}
