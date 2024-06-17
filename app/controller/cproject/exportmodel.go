package cproject

import (
	"fmt"
	"net/http"

	"github.com/pkg/errors"

	"projectforge.dev/projectforge/app"
	"projectforge.dev/projectforge/app/controller"
	"projectforge.dev/projectforge/app/controller/cutil"
	"projectforge.dev/projectforge/app/lib/metamodel/model"
	"projectforge.dev/projectforge/app/project/export/files"
	"projectforge.dev/projectforge/app/util"
	"projectforge.dev/projectforge/views/vexport"
)

func ProjectExportModelDetail(w http.ResponseWriter, r *http.Request) {
	controller.Act("project.export.model.detail", w, r, func(as *app.State, ps *cutil.PageState) (string, error) {
		prj, mdl, args, err := exportLoadModel(r, as, ps.Logger)
		if err != nil {
			return "", err
		}
		fls, err := files.ModelAll(mdl, prj, args, util.StringDefaultLinebreak)
		if err != nil {
			ps.Logger.Warnf("unable to generate files for model [%s]", mdl.Name)
		}
		bc := []string{"projects", prj.Key, fmt.Sprintf("Export||/p/%s/export", prj.Key), mdl.Title()}
		ps.SetTitleAndData(fmt.Sprintf("[%s] %s", prj.Key, mdl.Name), mdl)
		return controller.Render(r, as, &vexport.ModelDetail{Project: prj, Model: mdl, Files: fls}, ps, bc...)
	})
}

func ProjectExportModelSeedData(w http.ResponseWriter, r *http.Request) {
	controller.Act("project.export.seed.data", w, r, func(as *app.State, ps *cutil.PageState) (string, error) {
		prj, mdl, _, err := exportLoadModel(r, as, ps.Logger)
		if err != nil {
			return "", err
		}
		xbc := fmt.Sprintf("Export||/p/%s/export", prj.Key)
		mbc := fmt.Sprintf("%s||/p/%s/export/models/%s", mdl.Title(), prj.Key, mdl.Name)
		bc := []string{"projects", prj.Key, xbc, mbc, "Seed Data"}
		ps.SetTitleAndData(fmt.Sprintf("[%s] %s", prj.Key, mdl.Name), mdl)
		return controller.Render(r, as, &vexport.ModelSeedData{Project: prj, Model: mdl}, ps, bc...)
	})
}

func ProjectExportModelNew(w http.ResponseWriter, r *http.Request) {
	controller.Act("project.export.model.new", w, r, func(as *app.State, ps *cutil.PageState) (string, error) {
		prj, err := getProject(r, as)
		if err != nil {
			return "", err
		}
		mdl := &model.Model{}
		bc := []string{"projects", prj.Key, fmt.Sprintf("Export||/p/%s/export", prj.Key), "New"}
		ps.SetTitleAndData(fmt.Sprintf("[%s] New Model", prj.Key), mdl)
		return controller.Render(r, as, &vexport.ModelForm{Project: prj, Model: mdl, Examples: model.Examples}, ps, bc...)
	})
}

func ProjectExportModelCreate(w http.ResponseWriter, r *http.Request) {
	controller.Act("project.export.model.create", w, r, func(as *app.State, ps *cutil.PageState) (string, error) {
		prj, err := getProject(r, as)
		if err != nil {
			return "", err
		}

		frm, err := cutil.ParseForm(r, ps.RequestBody)
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
		return controller.FlashAndRedir(true, msg, u, ps)
	})
}

func ProjectExportModelForm(w http.ResponseWriter, r *http.Request) {
	controller.Act("project.export.model.form", w, r, func(as *app.State, ps *cutil.PageState) (string, error) {
		prj, mdl, _, err := exportLoadModel(r, as, ps.Logger)
		if err != nil {
			return "", err
		}
		bc := []string{
			"projects",
			prj.Key,
			fmt.Sprintf("Export||/p/%s/export", prj.Key),
			mdl.Title() + dblpipe + prj.WebPath() + "/export/models/" + mdl.Name,
			"Edit",
		}
		ps.SetTitleAndData(fmt.Sprintf("[%s] %s", prj.Key, mdl.Name), mdl)
		return controller.Render(r, as, &vexport.ModelForm{Project: prj, Model: mdl, Examples: model.Examples}, ps, bc...)
	})
}

func ProjectExportModelSave(w http.ResponseWriter, r *http.Request) {
	controller.Act("project.export.model.save", w, r, func(as *app.State, ps *cutil.PageState) (string, error) {
		prj, mdl, _, err := exportLoadModel(r, as, ps.Logger)
		if err != nil {
			return "", err
		}

		frm, err := cutil.ParseForm(r, ps.RequestBody)
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
		return controller.FlashAndRedir(true, msg, u, ps)
	})
}

func ProjectExportModelDelete(w http.ResponseWriter, r *http.Request) {
	controller.Act("project.export.model.delete", w, r, func(as *app.State, ps *cutil.PageState) (string, error) {
		prj, mdl, _, err := exportLoadModel(r, as, ps.Logger)
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
		return controller.FlashAndRedir(true, msg, fmt.Sprintf("/p/%s/export", prj.Key), ps)
	})
}
