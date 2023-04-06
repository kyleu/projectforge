package cproject

import (
	"fmt"

	"github.com/pkg/errors"
	"github.com/valyala/fasthttp"

	"projectforge.dev/projectforge/app"
	"projectforge.dev/projectforge/app/controller"
	"projectforge.dev/projectforge/app/controller/cutil"
	"projectforge.dev/projectforge/app/project/export/data"
	"projectforge.dev/projectforge/app/project/export/model"
	"projectforge.dev/projectforge/views/vexport"
)

func ProjectExportModelDeriveForm(rc *fasthttp.RequestCtx) {
	controller.Act("project.export.model.new", rc, func(as *app.State, ps *cutil.PageState) (string, error) {
		prj, err := getProject(rc, as)
		if err != nil {
			return "", err
		}

		mdl := &model.Model{}

		ps.Data = mdl

		bc := []string{"projects", prj.Key, fmt.Sprintf("Export||/p/%s/export", prj.Key), "Derive"}
		ps.Title = fmt.Sprintf("[%s] Derive Model", prj.Key)
		return controller.Render(rc, as, &vexport.DeriveForm{Project: prj}, ps, bc...)
	})
}

func ProjectExportModelDerive(rc *fasthttp.RequestCtx) {
	controller.Act("project.export.model.create", rc, func(as *app.State, ps *cutil.PageState) (string, error) {
		prj, err := getProject(rc, as)
		if err != nil {
			return "", err
		}

		frm, err := cutil.ParseForm(rc)
		if err != nil {
			return "", err
		}

		name := frm.GetStringOpt("name")
		pkg := frm.GetStringOpt("pkg")
		content := frm.GetStringOpt("content")
		if content == "" {
			return "", errors.New("content may not be empty")
		}

		mdl, err := data.Derive(name, pkg, content)
		if err != nil {
			return "", err
		}

		if cutil.QueryStringBool(rc, "save") {
			err = as.Services.Projects.SaveExportModel(as.Services.Projects.GetFilesystem(prj), mdl)
			if err != nil {
				return "", err
			}
		}

		msg := "model created successfully from input"
		u := fmt.Sprintf("/p/%s/export/models/%s", prj.Key, mdl.Name)
		return controller.FlashAndRedir(true, msg, u, rc, ps)
	})
}
