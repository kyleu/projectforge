package clib

import (
	"github.com/pkg/errors"
	"github.com/valyala/fasthttp"

	"{{{ .Package }}}/app"
	"{{{ .Package }}}/app/controller"
	"{{{ .Package }}}/app/controller/cutil"
	"{{{ .Package }}}/app/lib/notebook"
	"{{{ .Package }}}/views/vnotebook"
)

var notebookSvc *notebook.Service

func Notebook(rc *fasthttp.RequestCtx) {
	controller.Act("notebook", rc, func(as *app.State, ps *cutil.PageState) (string, error) {
		if notebookSvc == nil {
			notebookSvc = notebook.NewService()
		}
		status := notebookSvc.Status()
		if status == "running" {
			ps.SetTitleAndData("Notebook", "view-in-browser")
			return controller.Render(rc, as, &vnotebook.Notebook{}, ps, "notebook")
		}
		ps.SetTitleAndData("Notebook Options", status)
		return controller.Render(rc, as, &vnotebook.Options{}, ps, "notebook", "Options")
	})
}

func NotebookAction(rc *fasthttp.RequestCtx) {
	controller.Act("notebook.action", rc, func(as *app.State, ps *cutil.PageState) (string, error) {
		act, err := cutil.RCRequiredString(rc, "act", false)
		if err != nil {
			return "", err
		}
		switch act {
		case "start":
			err = notebookSvc.Start(as.Services.Exec)
			return controller.FlashAndRedir(true, "Notebook started", "/notebook", rc, ps)
		default:
			return "", errors.Errorf("invalid notebook action [%s]", act)
		}
	})
}
