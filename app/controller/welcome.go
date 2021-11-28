package controller

import (
	"strconv"

	"github.com/kyleu/projectforge/app"
	"github.com/kyleu/projectforge/app/controller/cutil"
	"github.com/kyleu/projectforge/app/project"
	"github.com/kyleu/projectforge/app/util"
	"github.com/kyleu/projectforge/views/vwelcome"
	"github.com/pkg/errors"
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
		ret, err := cutil.ParseForm(rc)
		if err != nil {
			return "", errors.Wrap(err, "unable to parse form")
		}

		prj := as.Services.Projects.ByPath(".")

		k := ret.GetStringOpt("key")
		if k != "" {
			prj.Key = k
		}
		prj.Name = ret.GetStringOpt("name")
		prj.Version = ret.GetStringOpt("version")
		if prj.Info == nil {
			prj.Info = &project.Info{}
		}
		prj.Info.Org = ret.GetStringOpt("org")
		prj.Package = ret.GetStringOpt("package")
		prj.Info.Homepage = ret.GetStringOpt("homepage")
		prj.Info.Sourcecode = ret.GetStringOpt("sourcecode")
		prj.Info.Summary = ret.GetStringOpt("summary")
		prj.Port, _ = strconv.Atoi(ret.GetStringOpt("port"))
		prj.Info.License = ret.GetStringOpt("license")

		err = as.Services.Projects.Save(prj)
		if err != nil {
			return "", errors.Wrap(err, "unable to save initial project")
		}

		return flashAndRedir(true, "project initialized", "/", rc, ps)
	})
}
