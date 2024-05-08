package cproject

import (
	"fmt"
	"net/http"

	"github.com/pkg/errors"

	"projectforge.dev/projectforge/app"
	"projectforge.dev/projectforge/app/controller"
	"projectforge.dev/projectforge/app/controller/cutil"
	"projectforge.dev/projectforge/app/lib/metamodel/model"
	"projectforge.dev/projectforge/app/project/export/derive"
	"projectforge.dev/projectforge/views/vexport"
)

func ProjectExportDeriveForm(w http.ResponseWriter, r *http.Request) {
	controller.Act("project.export.derive.form", w, r, func(as *app.State, ps *cutil.PageState) (string, error) {
		prj, err := getProject(r, as)
		if err != nil {
			return "", err
		}

		bc := []string{"projects", prj.Key, fmt.Sprintf("Export||/p/%s/export", prj.Key), "Derive"}
		ps.SetTitleAndData(fmt.Sprintf("[%s] Derive Model", prj.Key), &model.Model{})
		return controller.Render(r, as, &vexport.DeriveForm{Project: prj}, ps, bc...)
	})
}

func ProjectExportDerive(w http.ResponseWriter, r *http.Request) {
	controller.Act("project.export.derive", w, r, func(as *app.State, ps *cutil.PageState) (string, error) {
		prj, err := getProject(r, as)
		if err != nil {
			return "", err
		}

		frm, err := cutil.ParseForm(r, ps.RequestBody)
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
		if cutil.QueryStringBool(r, "save") {
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
		ps.SetTitleAndData(fmt.Sprintf("[%s] Derive Model", prj.Key), res)
		bc := []string{"projects", prj.Key, fmt.Sprintf("Export||/p/%s/export", prj.Key), "Derive"}
		return controller.Render(r, as, &vexport.DeriveForm{Project: prj, Result: res, Form: frm}, ps, bc...)
	})
}
