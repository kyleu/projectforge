package cproject

import (
	"net/http"

	"github.com/pkg/errors"

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

func Welcome(w http.ResponseWriter, r *http.Request) {
	controller.Act("welcome", w, r, func(as *app.State, ps *cutil.PageState) (string, error) {
		override := r.URL.Query().Get("override") == util.BoolTrue
		showLoad := override || r.URL.Query().Get("loaded") == util.BoolTrue
		if !showLoad {
			ps.HideMenu = true
			page := &vpage.Load{URL: "/welcome?loaded=true", Title: "Starting " + util.AppName, Message: "Checking some things..."}
			return controller.Render(w, r, as, page, ps, "Welcome")
		}
		ch := checks.CheckAll(ps.Context, as.Services.Modules.Modules().Keys(), ps.Logger, checks.Core(false).Keys()...).Errors()
		if len(ch) > 0 && (!override) {
			ps.SetTitleAndData("Missing Dependencies", ch.ErrorSummary())
			ps.HideMenu = true
			return controller.Render(w, r, as, &vwelcome.DepError{Results: ch}, ps, "Welcome||/welcome")
		}

		ps.SetTitleAndData("Welcome to "+util.AppName, welcomeMessage)
		ps.HideMenu = true
		return controller.Render(w, r, as, &vwelcome.Welcome{Project: as.Services.Projects.Default()}, ps, "Welcome||/welcome")
	})
}

var (
	activeWelcomeProject *project.Project
	activeWelcomeForm    util.ValueMap
)

func WelcomeLoad(w http.ResponseWriter, r *http.Request) {
	controller.Act("welcome.load", w, r, func(as *app.State, ps *cutil.PageState) (string, error) {
		frm, err := cutil.ParseForm(r, ps.RequestBody)
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
		r.URL.Query().Set("hasloaded", util.BoolTrue)
		page := &vpage.Load{URL: "/welcome/process", Title: "Generating and building your project"}
		return controller.Render(w, r, as, page, ps, "Welcome")
	})
}

func WelcomeProcess(w http.ResponseWriter, r *http.Request) {
	controller.Act("welcome.process", w, r, func(as *app.State, ps *cutil.PageState) (string, error) {
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

		return controller.FlashAndRedir(true, "project initialized", prj.WebPath(), w, ps)
	})
}
