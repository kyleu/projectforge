package cproject

import (
	"github.com/pkg/errors"
	"github.com/valyala/fasthttp"
	"projectforge.dev/projectforge/app/doctor/checks"
	"projectforge.dev/projectforge/views/vpage"

	"projectforge.dev/projectforge/app"
	"projectforge.dev/projectforge/app/controller"
	"projectforge.dev/projectforge/app/controller/cutil"
	"projectforge.dev/projectforge/app/project/action"
	"projectforge.dev/projectforge/app/util"
	"projectforge.dev/projectforge/views/vwelcome"
)

const welcomeMessage = "Welcome to " + util.AppName + "! View this page in a browser to get started."

func Welcome(rc *fasthttp.RequestCtx) {
	controller.Act("welcome", rc, func(as *app.State, ps *cutil.PageState) (string, error) {
		ch := checks.CheckAll(ps.Context, as.Services.Modules.Modules().Keys(), ps.Logger, "project", "projectforge", "repo").Errors()
		if len(ch) > 0 {
			ps.Title = util.AppName + ": Missing dependencies"
			ps.Data = ch.ErrorSummary()
			ps.HideMenu = true
			return controller.Render(rc, as, &vwelcome.DepError{Results: ch}, ps, "Welcome")
		}

		ps.Title = "Welcome to " + util.AppName
		ps.Data = welcomeMessage
		ps.HideMenu = true
		return controller.Render(rc, as, &vwelcome.Welcome{Project: as.Services.Projects.ByPath(".")}, ps, "Welcome")
	})
}

func WelcomeResult(rc *fasthttp.RequestCtx) {
	controller.Act("welcome.result", rc, func(as *app.State, ps *cutil.PageState) (string, error) {
		frm, err := cutil.ParseForm(rc)
		if err != nil {
			return "", err
		}

		if frm.GetStringOpt("hasloaded") != "true" {
			rc.URI().QueryArgs().Set("hasloaded", "true")
			page := &vpage.Load{URL: rc.URI().String(), Title: "Generating and building your project"}
			return controller.Render(rc, as, page, ps, "Welcome")
		}

		prj := as.Services.Projects.ByPath(".")

		mods, err := as.Services.Modules.GetModules(util.StringSplitAndTrim(frm.GetStringOpt("modules"), "||")...)
		if err != nil {
			return "", errors.Wrap(err, "can't parse modules")
		}
		prj.Modules = mods.Keys()
		frm["modules"] = mods.Keys()

		ret := action.Apply(ps.Context, actionParams(prj.Key, action.TypeCreate, frm, as, ps.Logger))
		if ret.HasErrors() {
			return "", errors.Wrap(ret.AsError(), "unable to build initial project")
		}

		return controller.FlashAndRedir(true, "project initialized", prj.WebPath(), rc, ps)
	})
}
