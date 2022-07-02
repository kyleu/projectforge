package cproject

import (
	"fmt"

	"github.com/pkg/errors"
	"github.com/valyala/fasthttp"

	"projectforge.dev/projectforge/app"
	"projectforge.dev/projectforge/app/controller"
	"projectforge.dev/projectforge/app/controller/cutil"
	"projectforge.dev/projectforge/app/export/files"
	"projectforge.dev/projectforge/app/export/model"
	"projectforge.dev/projectforge/views/vexport"
)

func ProjectExportModelDetail(rc *fasthttp.RequestCtx) {
	controller.Act("project.export.model.detail", rc, func(as *app.State, ps *cutil.PageState) (string, error) {
		prj, mdl, args, err := exportLoad(rc, as, ps.Logger)
		if err != nil {
			return "", err
		}
		ps.Data = mdl

		fls, err := files.ModelAll(mdl, args, true)
		if err != nil {
			ps.Logger.Warnf("unable to generate files for model [%s]", mdl.Name)
		}

		bc := []string{"projects", prj.Key, fmt.Sprintf("Export||/p/%s/export", prj.Key), mdl.Title()}
		ps.Title = fmt.Sprintf("[%s] %s", prj.Key, mdl.Name)
		return controller.Render(rc, as, &vexport.ModelDetail{Project: prj, Model: mdl, Files: fls}, ps, bc...)
	})
}

func ProjectExportModelSeedData(rc *fasthttp.RequestCtx) {
	controller.Act("project.export.seed.data", rc, func(as *app.State, ps *cutil.PageState) (string, error) {
		prj, mdl, _, err := exportLoad(rc, as, ps.Logger)
		if err != nil {
			return "", err
		}
		ps.Data = mdl

		xbc := fmt.Sprintf("Export||/p/%s/export", prj.Key)
		mbc := fmt.Sprintf("%s||/p/%s/export/models/%s", mdl.Title(), prj.Key, mdl.Name)
		bc := []string{"projects", prj.Key, xbc, mbc, "Seed Data"}
		ps.Title = fmt.Sprintf("[%s] %s", prj.Key, mdl.Name)
		return controller.Render(rc, as, &vexport.ModelSeedData{Project: prj, Model: mdl}, ps, bc...)
	})
}

func ProjectExportModelNew(rc *fasthttp.RequestCtx) {
	controller.Act("project.export.model.new", rc, func(as *app.State, ps *cutil.PageState) (string, error) {
		prj, err := getProject(rc, as)
		if err != nil {
			return "", err
		}

		mdl := &model.Model{}

		ps.Data = mdl

		bc := []string{"projects", prj.Key, fmt.Sprintf("Export||/p/%s/export", prj.Key), "New"}
		ps.Title = fmt.Sprintf("[%s] New Model", prj.Key)
		return controller.Render(rc, as, &vexport.ModelForm{Project: prj, Model: mdl, Examples: model.Examples}, ps, bc...)
	})
}

func ProjectExportModelCreate(rc *fasthttp.RequestCtx) {
	controller.Act("project.export.model.create", rc, func(as *app.State, ps *cutil.PageState) (string, error) {
		prj, err := getProject(rc, as)
		if err != nil {
			return "", err
		}

		frm, err := cutil.ParseForm(rc)
		if err != nil {
			return "", err
		}

		mdl := &model.Model{}
		err = exportModelFromForm(frm, mdl)
		if err != nil {
			return "", errors.Wrap(err, "unable to parse model from form")
		}

		err = as.Services.Projects.SaveExportModel(as.Services.Projects.GetFilesystem(prj), mdl, ps.Logger)
		if err != nil {
			return "", err
		}

		msg := "model created successfully"
		u := fmt.Sprintf("/p/%s/export/models/%s", prj.Key, mdl.Name)
		return controller.FlashAndRedir(true, msg, u, rc, ps)
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

		msg := "model saved successfully"
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

		msg := "model deleted successfully"
		return controller.FlashAndRedir(true, msg, fmt.Sprintf("/p/%s/export", prj.Key), rc, ps)
	})
}
