package cproject

import (
	"fmt"

	"github.com/pkg/errors"
	"github.com/valyala/fasthttp"

	"projectforge.dev/projectforge/app"
	"projectforge.dev/projectforge/app/controller"
	"projectforge.dev/projectforge/app/controller/cutil"
	"projectforge.dev/projectforge/app/project/export/files"
	"projectforge.dev/projectforge/app/project/export/model"
	"projectforge.dev/projectforge/app/util"
	"projectforge.dev/projectforge/views/vexport"
)

func ProjectExportModelDetail(rc *fasthttp.RequestCtx) {
	controller.Act("project.export.model.detail", rc, func(as *app.State, ps *cutil.PageState) (string, error) {
		prj, mdl, args, err := exportLoadModel(rc, as, ps.Logger)
		if err != nil {
			return "", err
		}
		fls, err := files.ModelAll(mdl, prj, args, true, util.StringDefaultLinebreak)
		if err != nil {
			ps.Logger.Warnf("unable to generate files for model [%s]", mdl.Name)
		}
		bc := []string{"projects", prj.Key, fmt.Sprintf("Export||/p/%s/export", prj.Key), mdl.Title()}
		ps.Title = fmt.Sprintf("[%s] %s", prj.Key, mdl.Name)
		ps.Data = mdl
		return controller.Render(rc, as, &vexport.ModelDetail{Project: prj, Model: mdl, Files: fls}, ps, bc...)
	})
}

func ProjectExportModelSeedData(rc *fasthttp.RequestCtx) {
	controller.Act("project.export.seed.data", rc, func(as *app.State, ps *cutil.PageState) (string, error) {
		prj, mdl, _, err := exportLoadModel(rc, as, ps.Logger)
		if err != nil {
			return "", err
		}
		xbc := fmt.Sprintf("Export||/p/%s/export", prj.Key)
		mbc := fmt.Sprintf("%s||/p/%s/export/models/%s", mdl.Title(), prj.Key, mdl.Name)
		bc := []string{"projects", prj.Key, xbc, mbc, "Seed Data"}
		ps.Title = fmt.Sprintf("[%s] %s", prj.Key, mdl.Name)
		ps.Data = mdl
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
		bc := []string{"projects", prj.Key, fmt.Sprintf("Export||/p/%s/export", prj.Key), "New"}
		ps.Title = fmt.Sprintf("[%s] New Model", prj.Key)
		ps.Data = mdl
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

		pfs, err := as.Services.Projects.GetFilesystem(prj)
		if err != nil {
			return "", err
		}
		err = as.Services.Projects.SaveExportModel(pfs, mdl)
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
		prj, mdl, _, err := exportLoadModel(rc, as, ps.Logger)
		if err != nil {
			return "", err
		}
		bc := []string{
			"projects",
			prj.Key,
			fmt.Sprintf("Export||/p/%s/export", prj.Key),
			mdl.Title() + "||/p/" + prj.Key + "/export/models/" + mdl.Name,
			"Edit",
		}
		ps.Title = fmt.Sprintf("[%s] %s", prj.Key, mdl.Name)
		ps.Data = mdl
		return controller.Render(rc, as, &vexport.ModelForm{Project: prj, Model: mdl, Examples: model.Examples}, ps, bc...)
	})
}

func ProjectExportModelSave(rc *fasthttp.RequestCtx) {
	controller.Act("project.export.model.save", rc, func(as *app.State, ps *cutil.PageState) (string, error) {
		prj, mdl, _, err := exportLoadModel(rc, as, ps.Logger)
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

		pfs, err := as.Services.Projects.GetFilesystem(prj)
		if err != nil {
			return "", err
		}
		err = as.Services.Projects.SaveExportModel(pfs, mdl)
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
		prj, mdl, _, err := exportLoadModel(rc, as, ps.Logger)
		if err != nil {
			return "", err
		}

		pfs, err := as.Services.Projects.GetFilesystem(prj)
		if err != nil {
			return "", err
		}
		err = as.Services.Projects.DeleteExportModel(pfs, mdl.Name, ps.Logger)
		if err != nil {
			return "", err
		}

		msg := "model deleted successfully"
		return controller.FlashAndRedir(true, msg, fmt.Sprintf("/p/%s/export", prj.Key), rc, ps)
	})
}
