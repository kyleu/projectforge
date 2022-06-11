package cproject

import (
	"fmt"

	"github.com/pkg/errors"
	"github.com/valyala/fasthttp"

	"projectforge.dev/projectforge/app"
	"projectforge.dev/projectforge/app/controller"
	"projectforge.dev/projectforge/app/controller/cutil"
	"projectforge.dev/projectforge/app/export/model"
	"projectforge.dev/projectforge/views/vexport"
)

func ProjectExportModelDetail(rc *fasthttp.RequestCtx) {
	controller.Act("project.export.model.detail", rc, func(as *app.State, ps *cutil.PageState) (string, error) {
		prj, mdl, _, err := exportLoad(rc, as, ps.Logger)
		if err != nil {
			return "", err
		}
		ps.Data = mdl

		bc := []string{"projects", prj.Key, fmt.Sprintf("Export||/p/%s/export", prj.Key), mdl.Title()}
		ps.Title = fmt.Sprintf("[%s] %s", prj.Key, mdl.Name)
		return controller.Render(rc, as, &vexport.ModelDetail{Project: prj, Model: mdl}, ps, bc...)
	})
}

func ProjectExportModelForm(rc *fasthttp.RequestCtx) {
	controller.Act("project.export.model.form", rc, func(as *app.State, ps *cutil.PageState) (string, error) {
		prj, mdl, _, err := exportLoad(rc, as, ps.Logger)
		if err != nil {
			return "", err
		}

		ps.Data = mdl

		bc := []string{
			"projects",
			prj.Key,
			fmt.Sprintf("Export||/p/%s/export", prj.Key),
			mdl.Title() + "||/p/" + prj.Key + "/export/models/" + mdl.Name,
			"Edit",
		}
		ps.Title = fmt.Sprintf("[%s] %s", prj.Key, mdl.Name)
		return controller.Render(rc, as, &vexport.ModelForm{Project: prj, Model: mdl, Examples: model.Examples}, ps, bc...)
	})
}

func ProjectExportModelSave(rc *fasthttp.RequestCtx) {
	controller.Act("project.export.model.save", rc, func(as *app.State, ps *cutil.PageState) (string, error) {
		prj, mdl, _, err := exportLoad(rc, as, ps.Logger)
		if err != nil {
			return "", err
		}

		frm, err := cutil.ParseForm(rc)
		if err != nil {
			return "", err
		}

		err = exportModelFromForm(frm, mdl)
		if err != nil {
			return "", errors.Wrap(err, "unable to parse model from form")
		}

		err = as.Services.Projects.SaveExportModel(as.Services.Projects.GetFilesystem(prj), mdl, ps.Logger)
		if err != nil {
			return "", err
		}

		msg := fmt.Sprintf("model saved successfully")
		u := fmt.Sprintf("/p/%s/export/models/%s", prj.Key, mdl.Name)
		return controller.FlashAndRedir(true, msg, u, rc, ps)
	})
}

func ProjectExportModelDelete(rc *fasthttp.RequestCtx) {
	controller.Act("project.export.model.delete", rc, func(as *app.State, ps *cutil.PageState) (string, error) {
		prj, mdl, _, err := exportLoad(rc, as, ps.Logger)
		if err != nil {
			return "", err
		}

		err = as.Services.Projects.DeleteExportModel(as.Services.Projects.GetFilesystem(prj), mdl.Name, ps.Logger)
		if err != nil {
			return "", err
		}

		msg := fmt.Sprintf("model saved successfully")
		return controller.FlashAndRedir(true, msg, fmt.Sprintf("/p/%s/export", prj.Key), rc, ps)
	})
}
