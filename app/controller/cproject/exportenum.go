package cproject

import (
	"fmt"

	"github.com/pkg/errors"
	"github.com/valyala/fasthttp"

	"projectforge.dev/projectforge/app"
	"projectforge.dev/projectforge/app/controller"
	"projectforge.dev/projectforge/app/controller/cutil"
	"projectforge.dev/projectforge/app/project/export/enum"
	"projectforge.dev/projectforge/app/project/export/files/goenum"
	"projectforge.dev/projectforge/app/util"
	"projectforge.dev/projectforge/views/vexport"
)

func ProjectExportEnumDetail(rc *fasthttp.RequestCtx) {
	controller.Act("project.export.enum.detail", rc, func(as *app.State, ps *cutil.PageState) (string, error) {
		prj, e, _, err := exportLoadEnum(rc, as, ps.Logger)
		if err != nil {
			return "", err
		}
		fl, err := goenum.Enum(e, true, util.StringDefaultLinebreak)
		if err != nil {
			ps.Logger.Warnf("unable to generate files for enum [%s]", e.Name)
		}

		bc := []string{"projects", prj.Key, fmt.Sprintf("Export||/p/%s/export", prj.Key), e.Title()}
		ps.Title = fmt.Sprintf("[%s] %s", prj.Key, e.Name)
		ps.Data = e
		return controller.Render(rc, as, &vexport.EnumDetail{Project: prj, Enum: e, File: fl}, ps, bc...)
	})
}

func ProjectExportEnumNew(rc *fasthttp.RequestCtx) {
	controller.Act("project.export.enum.new", rc, func(as *app.State, ps *cutil.PageState) (string, error) {
		prj, err := getProject(rc, as)
		if err != nil {
			return "", err
		}

		e := &enum.Enum{}
		bc := []string{"projects", prj.Key, fmt.Sprintf("Export||/p/%s/export", prj.Key), "New"}
		ps.Title = fmt.Sprintf("[%s] New Enum", prj.Key)
		ps.Data = e
		return controller.Render(rc, as, &vexport.EnumForm{Project: prj, Enum: e}, ps, bc...)
	})
}

func ProjectExportEnumCreate(rc *fasthttp.RequestCtx) {
	controller.Act("project.export.enum.create", rc, func(as *app.State, ps *cutil.PageState) (string, error) {
		prj, err := getProject(rc, as)
		if err != nil {
			return "", err
		}

		frm, err := cutil.ParseForm(rc)
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
		return controller.FlashAndRedir(true, msg, u, rc, ps)
	})
}

func ProjectExportEnumForm(rc *fasthttp.RequestCtx) {
	controller.Act("project.export.enum.form", rc, func(as *app.State, ps *cutil.PageState) (string, error) {
		prj, e, _, err := exportLoadEnum(rc, as, ps.Logger)
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
		ps.Title = fmt.Sprintf("[%s] %s", prj.Key, e.Name)
		ps.Data = e
		return controller.Render(rc, as, &vexport.EnumForm{Project: prj, Enum: e}, ps, bc...)
	})
}

func ProjectExportEnumSave(rc *fasthttp.RequestCtx) {
	controller.Act("project.export.enum.save", rc, func(as *app.State, ps *cutil.PageState) (string, error) {
		prj, e, _, err := exportLoadEnum(rc, as, ps.Logger)
		if err != nil {
			return "", err
		}

		frm, err := cutil.ParseForm(rc)
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
		return controller.FlashAndRedir(true, msg, u, rc, ps)
	})
}

func ProjectExportEnumDelete(rc *fasthttp.RequestCtx) {
	controller.Act("project.export.enum.delete", rc, func(as *app.State, ps *cutil.PageState) (string, error) {
		prj, mdl, _, err := exportLoadEnum(rc, as, ps.Logger)
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
		return controller.FlashAndRedir(true, msg, fmt.Sprintf("/p/%s/export", prj.Key), rc, ps)
	})
}
