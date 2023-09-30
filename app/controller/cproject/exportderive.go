package cproject

import (
	"fmt"

	"github.com/pkg/errors"
	"github.com/valyala/fasthttp"

	"projectforge.dev/projectforge/app"
	"projectforge.dev/projectforge/app/controller"
	"projectforge.dev/projectforge/app/controller/cutil"
	"projectforge.dev/projectforge/app/project/export/derive"
	"projectforge.dev/projectforge/app/project/export/model"
	"projectforge.dev/projectforge/views/vexport"
)

func ProjectExportDeriveForm(rc *fasthttp.RequestCtx) {
	controller.Act("project.export.derive.form", rc, func(as *app.State, ps *cutil.PageState) (string, error) {
		prj, err := getProject(rc, as)
		if err != nil {
			return "", err
		}

		bc := []string{"projects", prj.Key, fmt.Sprintf("Export||/p/%s/export", prj.Key), "Derive"}
		ps.Title = fmt.Sprintf("[%s] Derive Model", prj.Key)
		ps.Data = &model.Model{}
		return controller.Render(rc, as, &vexport.DeriveForm{Project: prj}, ps, bc...)
	})
}

func ProjectExportDerive(rc *fasthttp.RequestCtx) {
	controller.Act("project.export.derive", rc, func(as *app.State, ps *cutil.PageState) (string, error) {
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

		res := derive.Derive(name, pkg, content, ps.Logger)
		if cutil.QueryStringBool(rc, "save") {
			pfs, err := as.Services.Projects.GetFilesystem(prj)
			if err != nil {
				return "", err
			}
			for _, s := range res {
				if len(s.Models) > 0 {
					for _, mdl := range s.Models {
						err = as.Services.Projects.SaveExportModel(pfs, mdl)
						if err != nil {
							return "", err
						}
					}
				}
			}
		}
		ps.Title = fmt.Sprintf("[%s] Derive Model", prj.Key)
		ps.Data = res
		bc := []string{"projects", prj.Key, fmt.Sprintf("Export||/p/%s/export", prj.Key), "Derive"}
		return controller.Render(rc, as, &vexport.DeriveForm{Project: prj, Result: res, Form: frm}, ps, bc...)
	})
}
