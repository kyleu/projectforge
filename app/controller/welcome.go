package controller

import (
	"github.com/kyleu/projectforge/app"
	"github.com/kyleu/projectforge/app/controller/cutil"
	"github.com/kyleu/projectforge/app/util"
	"github.com/kyleu/projectforge/views/vwelcome"
	"github.com/valyala/fasthttp"
)

const welcomeMessage = "Welcome to " + util.AppName + "! View this page in a browser to get started."

func Welcome(rc *fasthttp.RequestCtx) {
	act("welcome", rc, func(as *app.State, ps *cutil.PageState) (string, error) {
		ps.Title = "Welcome to " + util.AppName
		ps.Data = welcomeMessage
		ps.HideMenu = true
		return render(rc, as, &vwelcome.Welcome{Project: as.Services.Projects.ByPath(".")}, ps, "Welcome")
	})
}

func WelcomeResult(rc *fasthttp.RequestCtx) {
	act("welcome.result", rc, func(as *app.State, ps *cutil.PageState) (string, error) {
		return flashAndRedir(true, "project initialized", "/", rc, ps)
	})
}
