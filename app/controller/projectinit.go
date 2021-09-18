package controller

import (
	"fmt"
	"os"

	"github.com/kyleu/projectforge/app/util"
	"github.com/kyleu/projectforge/views/vinit"
	"github.com/valyala/fasthttp"

	"github.com/kyleu/projectforge/app/controller/cutil"

	"github.com/kyleu/projectforge/app"
)

func ProjectInitWarning(rc *fasthttp.RequestCtx) {
	act("project.init.warning", rc, func(as *app.State, ps *cutil.PageState) (string, error) {
		wd, _ := os.Getwd()
		ps.Title = "Initialize"
		msg := "%s has not been initialized in this directory. To initialize [%s], POST to this URL"
		ps.Data = fmt.Sprintf(msg, util.AppName, wd)
		return render(rc, as, &vinit.Form{Dir: wd}, ps, "Initialize")
	})
}

func ProjectInit(rc *fasthttp.RequestCtx) {
	act("project.init", rc, func(as *app.State, ps *cutil.PageState) (string, error) {
		return "/", as.Services.Projects.Init()
	})
}
