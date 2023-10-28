package cproject

import (
	"github.com/pkg/errors"
	"github.com/valyala/fasthttp"

	"projectforge.dev/projectforge/app"
	"projectforge.dev/projectforge/app/controller"
	"projectforge.dev/projectforge/app/controller/cutil"
	"projectforge.dev/projectforge/app/doctor/checks"
	"projectforge.dev/projectforge/app/project"
	"projectforge.dev/projectforge/app/project/action"
	"projectforge.dev/projectforge/app/util"
	"projectforge.dev/projectforge/views/vpage"
	"projectforge.dev/projectforge/views/vwelcome"
)

const welcomeMessage = "Welcome to " + util.AppName + "! View this page in a browser to get started."

func Welcome(rc *fasthttp.RequestCtx) {
	controller.Act("welcome", rc, func(as *app.State, ps *cutil.PageState) (string, error) {
		override := string(rc.URI().QueryArgs().Peek("override")) == util.BoolTrue
		showLoad := override || string(rc.URI().QueryArgs().Peek("loaded")) == util.BoolTrue
		if !showLoad {
			ps.HideMenu = true
			page := &vpage.Load{URL: "/welcome?loaded=true", Title: "Starting " + util.AppName, Message: "Checking some things..."}
			return controller.Render(rc, as, page, ps, "Welcome")
		}
		ch := checks.CheckAll(ps.Context, as.Services.Modules.Modules().Keys(), ps.Logger, checks.Core(false).Keys()...).Errors()
		if len(ch) > 0 && (!override) {
			ps.SetTitleAndData("Missing Dependencies", ch.ErrorSummary())
			ps.HideMenu = true
			return controller.Render(rc, as, &vwelcome.DepError{Results: ch}, ps, "Welcome||/welcome")
		}

		ps.SetTitleAndData("Welcome to "+util.AppName, welcomeMessage)
		ps.HideMenu = true
		return controller.Render(rc, as, &vwelcome.Welcome{Project: as.Services.Projects.Default()}, ps, "Welcome||/welcome")
	})
}

var (
	activeWelcomeProject *project.Project
	activeWelcomeForm    util.ValueMap
)

func WelcomeLoad(rc *fasthttp.RequestCtx) {
	controller.Act("welcome.load", rc, func(as *app.State, ps *cutil.PageState) (string, error) {
		frm, err := cutil.ParseForm(rc)
		if err != nil {
			return "", err
		}

		prj := as.Services.Projects.Default()
		mods, err := as.Services.Modules.GetModules(util.StringSplitAndTrim(frm.GetStringOpt("modules"), "||")...)
		if err != nil {
			return "", errors.Wrap(err, "can't parse modules")
		}
		prj.Modules = mods.Keys()
		frm["modules"] = mods.Keys()

		activeWelcomeProject = prj
		activeWelcomeForm = frm
		rc.URI().QueryArgs().Set("hasloaded", util.BoolTrue)
		page := &vpage.Load{URL: "/welcome/process", Title: "Generating and building your project"}
		return controller.Render(rc, as, page, ps, "Welcome")
	})
}

func WelcomeProcess(rc *fasthttp.RequestCtx) {
	controller.Act("welcome.process", rc, func(as *app.State, ps *cutil.PageState) (string, error) {
		frm := activeWelcomeForm
		prj := activeWelcomeProject
		if prj == nil {
			return "", errors.New("Project loading hasn't occurred")
		}
		if frm == nil {
			return "", errors.New("Form loading hasn't occurred")
		}
		ret := action.Apply(ps.Context, actionParams(prj.Key, action.TypeCreate, frm, as, ps.Logger))
		if ret.HasErrors() {
			return "", errors.Wrap(ret.AsError(), "unable to build initial project")
		}
		activeWelcomeProject = nil
		activeWelcomeForm = nil

		return controller.FlashAndRedir(true, "project initialized", prj.WebPath(), rc, ps)
	})
}
