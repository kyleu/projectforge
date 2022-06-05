package controller

import (
	"fmt"

	"github.com/pkg/errors"
	"github.com/valyala/fasthttp"
	"projectforge.dev/projectforge/app"
	"projectforge.dev/projectforge/app/controller/cutil"
	"projectforge.dev/projectforge/views/vexport"
)

func ProjectExportOverview(rc *fasthttp.RequestCtx) {
	act("project.export.overview", rc, func(as *app.State, ps *cutil.PageState) (string, error) {
		prj, err := getProject(rc, as)
		if err != nil {
			return "", err
		}
		args, err := prj.Info.ModuleArgExport()
		if err != nil {
			return "", err
		}
		ps.Data = args

		bc := []string{"projects", prj.Key, "Export"}
		ps.Title = fmt.Sprintf("[%s] Export", prj.Key)
		return render(rc, as, &vexport.Overview{Project: prj, Args: args}, ps, bc...)
	})
}

func ProjectExportModelDetail(rc *fasthttp.RequestCtx) {
	act("project.export.model.detail", rc, func(as *app.State, ps *cutil.PageState) (string, error) {
		prj, mdl, _, err := exportLoad(rc, as)
		if err != nil {
			return "", err
		}
		ps.Data = mdl

		bc := []string{"projects", prj.Key, fmt.Sprintf("Export||/p/%s/export", prj.Key), mdl.Title()}
		ps.Title = fmt.Sprintf("[%s] %s", prj.Key, mdl.Name)
		return render(rc, as, &vexport.ModelDetail{Project: prj, Model: mdl}, ps, bc...)
	})
}

func ProjectExportModelForm(rc *fasthttp.RequestCtx) {
	act("project.export.model.form", rc, func(as *app.State, ps *cutil.PageState) (string, error) {
		prj, mdl, _, err := exportLoad(rc, as)
		if err != nil {
			return "", err
		}

		ps.Data = mdl

		bc := []string{"projects", prj.Key, fmt.Sprintf("Export||/p/%s/export", prj.Key), mdl.Title() + "||/p/" + prj.Key + "/export/" + mdl.Name, "Edit"}
		ps.Title = fmt.Sprintf("[%s] %s", prj.Key, mdl.Name)
		return render(rc, as, &vexport.ModelForm{Project: prj, Model: mdl}, ps, bc...)
	})
}

func ProjectExportModelSave(rc *fasthttp.RequestCtx) {
	act("project.export.model.save", rc, func(as *app.State, ps *cutil.PageState) (string, error) {
		prj, mdl, args, err := exportLoad(rc, as)
		if err != nil {
			return "", err
		}

		frm, err := cutil.ParseForm(rc)
		if err != nil {
			return "", errors.Wrap(err, "unable to parse form")
		}

		err = exportModelFromForm(frm, mdl)
		if err != nil {
			return "", errors.Wrap(err, "unable to parse model from form")
		}

		args.Models = args.Models.Replace(mdl)

		msg := fmt.Sprintf("model saved successfully")
		u := fmt.Sprintf("/p/%s/export/%s", prj.Key, mdl.Name)
		return flashAndRedir(true, msg, u, rc, ps)
	})
}
