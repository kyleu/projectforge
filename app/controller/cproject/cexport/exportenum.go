package cexport

import (
	"fmt"
	"net/http"

	"github.com/pkg/errors"

	"projectforge.dev/projectforge/app"
	"projectforge.dev/projectforge/app/controller"
	"projectforge.dev/projectforge/app/controller/cproject"
	"projectforge.dev/projectforge/app/controller/cutil"
	"projectforge.dev/projectforge/app/lib/metamodel/enum"
	"projectforge.dev/projectforge/app/project/export/files/goenum"
	"projectforge.dev/projectforge/app/util"
	"projectforge.dev/projectforge/views/vexport"
)

func ProjectExportEnumDetail(w http.ResponseWriter, r *http.Request) {
	controller.Act("project.export.enum.detail", w, r, func(as *app.State, ps *cutil.PageState) (string, error) {
		prj, e, err := exportLoadEnum(r, as, ps.Logger)
		if err != nil {
			return "", err
		}
		fl, err := goenum.Enum(e, util.StringDefaultLinebreak)
		if err != nil {
			ps.Logger.Warnf("unable to generate [%s] files for enum [%s]", prj.Key, e.Name)
		}

		bc := []string{"projects", prj.Key, fmt.Sprintf("Export||/p/%s/export", prj.Key), e.Title()}
		ps.SetTitleAndData(fmt.Sprintf("[%s] %s", prj.Key, e.Name), e)
		return controller.Render(r, as, &vexport.EnumDetail{BaseURL: prj.WebPathEnums(), Enum: e, File: fl}, ps, bc...)
	})
}

func ProjectExportEnumNew(w http.ResponseWriter, r *http.Request) {
	controller.Act("project.export.enum.new", w, r, func(as *app.State, ps *cutil.PageState) (string, error) {
		prj, err := cproject.GetProjectWithArgs(r, as, ps.Logger)
		if err != nil {
			return "", err
		}

		e := &enum.Enum{}
		bc := []string{"projects", prj.Key, fmt.Sprintf("Export||/p/%s/export", prj.Key), "New"}
		ps.SetTitleAndData(fmt.Sprintf("[%s] New Enum", prj.Key), e)
		return controller.Render(r, as, &vexport.EnumForm{BaseURL: prj.WebPathEnums(), Enum: e, Examples: enum.Examples}, ps, bc...)
	})
}

func ProjectExportEnumCreate(w http.ResponseWriter, r *http.Request) {
	controller.Act("project.export.enum.create", w, r, func(as *app.State, ps *cutil.PageState) (string, error) {
		prj, err := cproject.GetProjectWithArgs(r, as, ps.Logger)
		if err != nil {
			return "", err
		}

		frm, err := cutil.ParseForm(r, ps.RequestBody)
		if err != nil {
			return "", err
		}

		mdl := &enum.Enum{}
		err = exportEnumFromForm(frm, mdl)
		if err != nil {
			return "", errors.Wrap(err, "unable to parse enum from form")
		}

		pfs, err := as.Services.Projects.GetFilesystem(prj)
		if err != nil {
			return "", err
		}
		err = as.Services.Projects.SaveExportEnum(pfs, mdl)
		if err != nil {
			return "", err
		}

		msg := "enum created successfully"
		u := fmt.Sprintf("/p/%s/export/enums/%s", prj.Key, mdl.Name)
		return controller.FlashAndRedir(true, msg, u, ps)
	})
}

func ProjectExportEnumForm(w http.ResponseWriter, r *http.Request) {
	controller.Act("project.export.enum.form", w, r, func(as *app.State, ps *cutil.PageState) (string, error) {
		prj, e, err := exportLoadEnum(r, as, ps.Logger)
		if err != nil {
			return "", err
		}

		bc := []string{
			"projects",
			prj.Key,
			fmt.Sprintf("Export||/p/%s/export", prj.Key),
			fmt.Sprintf("%s||/p/%s/export/enums/%s", e.Title(), prj.Key, e.Name),
			"Edit",
		}
		ps.SetTitleAndData(fmt.Sprintf("[%s] %s", prj.Key, e.Name), e)
		return controller.Render(r, as, &vexport.EnumForm{BaseURL: prj.WebPathEnums(), Enum: e, Examples: enum.Examples}, ps, bc...)
	})
}

func ProjectExportEnumSave(w http.ResponseWriter, r *http.Request) {
	controller.Act("project.export.enum.save", w, r, func(as *app.State, ps *cutil.PageState) (string, error) {
		prj, e, err := exportLoadEnum(r, as, ps.Logger)
		if err != nil {
			return "", err
		}

		frm, err := cutil.ParseForm(r, ps.RequestBody)
		if err != nil {
			return "", err
		}

		err = exportEnumFromForm(frm, e)
		if err != nil {
			return "", errors.Wrap(err, "unable to parse enum from form")
		}
		pfs, err := as.Services.Projects.GetFilesystem(prj)
		if err != nil {
			return "", err
		}
		err = as.Services.Projects.SaveExportEnum(pfs, e)
		if err != nil {
			return "", err
		}

		msg := "enum saved successfully"
		u := fmt.Sprintf("/p/%s/export/enums/%s", prj.Key, e.Name)
		return controller.FlashAndRedir(true, msg, u, ps)
	})
}

func ProjectExportEnumDelete(w http.ResponseWriter, r *http.Request) {
	controller.Act("project.export.enum.delete", w, r, func(as *app.State, ps *cutil.PageState) (string, error) {
		prj, mdl, err := exportLoadEnum(r, as, ps.Logger)
		if err != nil {
			return "", err
		}
		pfs, err := as.Services.Projects.GetFilesystem(prj)
		if err != nil {
			return "", err
		}
		err = as.Services.Projects.DeleteExportEnum(pfs, mdl.Name, ps.Logger)
		if err != nil {
			return "", err
		}

		msg := "enum deleted successfully"
		return controller.FlashAndRedir(true, msg, fmt.Sprintf("/p/%s/export", prj.Key), ps)
	})
}
