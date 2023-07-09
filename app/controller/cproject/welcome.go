package cproject

import (
	"github.com/pkg/errors"
	"github.com/valyala/fasthttp"
	"projectforge.dev/projectforge/app/project"

	"projectforge.dev/projectforge/app"
	"projectforge.dev/projectforge/app/controller"
	"projectforge.dev/projectforge/app/controller/cutil"
	"projectforge.dev/projectforge/app/doctor/checks"
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
		ch := checks.CheckAll(ps.Context, as.Services.Modules.Modules().Keys(), ps.Logger, "project", "projectforge", "repo").Errors()
		if len(ch) > 0 && (!override) {
			ps.Title = "Missing Dependencies"
			ps.Data = ch.ErrorSummary()
			ps.HideMenu = true
			return controller.Render(rc, as, &vwelcome.DepError{Results: ch}, ps, "Welcome||/welcome")
		}

		ps.Title = "Welcome to " + util.AppName
		ps.Data = welcomeMessage
		ps.HideMenu = true
		return controller.Render(rc, as, &vwelcome.Welcome{Project: as.Services.Projects.Default()}, ps, "Welcome||/welcome")
	})
}

var activeWelcomeProject *project.Project
var activeWelcomeForm util.ValueMap

func WelcomeResult(rc *fasthttp.RequestCtx) {
	controller.Act("welcome.result", rc, func(as *app.State, ps *cutil.PageState) (string, error) {
		frm, err := cutil.ParseForm(rc)
		if err != nil {
			return "", err
		}

		if frm.GetStringOpt("hasloaded") != util.BoolTrue {
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
			page := &vpage.Load{URL: rc.URI().String(), Title: "Generating and building your project"}
			return controller.Render(rc, as, page, ps, "Welcome")
		}

		if activeWelcomeProject == nil {
			return "", errors.New("Project loading hasn't occurred")
		}
		if activeWelcomeForm == nil {
			return "", errors.New("Form loading hasn't occurred")
		}
		ret := action.Apply(ps.Context, actionParams(activeWelcomeProject.Key, action.TypeCreate, activeWelcomeForm, as, ps.Logger))
		if ret.HasErrors() {
			return "", errors.Wrap(ret.AsError(), "unable to build initial project")
		}
		activeWelcomeProject = nil
		activeWelcomeForm = nil

		return controller.FlashAndRedir(true, "project initialized", activeWelcomeProject.WebPath(), rc, ps)
	})
}
