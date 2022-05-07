package controller

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/pkg/errors"
	"github.com/valyala/fasthttp"

	"projectforge.dev/projectforge/app"
	"projectforge.dev/projectforge/app/lib/theme"
	"projectforge.dev/projectforge/app/project"
	"projectforge.dev/projectforge/app/util"
)

func projectFromForm(frm util.ValueMap, prj *project.Project) error {
	get := func(k string, def string) string {
		x := frm.GetStringOpt(k)
		if x == "" {
			return def
		}
		return x
	}
	prj.Name = get("name", prj.Name)
	prj.Version = get("version", prj.Version)
	prj.Package = get("package", prj.Package)
	if prj.Package == "" {
		prj.Package = "github.com/org/" + prj.Key
	}
	prj.Args = get("args", prj.Args)
	prj.Port, _ = strconv.Atoi(get("port", fmt.Sprintf("%d", prj.Port)))
	if prj.Port == 0 {
		prj.Port = 10000
	}
	prj.Modules = util.StringSplitAndTrim(get("modules", strings.Join(prj.Modules, "||")), "||")
	if len(prj.Modules) == 0 {
		prj.Modules = []string{"core"}
	}
	prj.Ignore = util.StringSplitAndTrim(get("ignore", strings.Join(prj.Ignore, ",")), ",")
	prj.Path = get("path", prj.Path)

	if prj.Info == nil {
		prj.Info = &project.Info{}
	}
	prj.Info.Org = get("org", prj.Info.Org)
	prj.Info.AuthorID = get("authorID", prj.Info.AuthorID)
	prj.Info.AuthorName = get("authorName", prj.Info.AuthorName)
	prj.Info.AuthorEmail = get("authorEmail", prj.Info.AuthorEmail)
	prj.Info.License = get("license", prj.Info.License)
	prj.Info.Homepage = get("homepage", prj.Info.Homepage)
	prj.Info.Sourcecode = get("sourcecode", prj.Info.Sourcecode)
	prj.Info.Summary = get("summary", prj.Info.Summary)
	prj.Info.Description = get("description", prj.Info.Description)
	prj.Info.CI = get("ci", prj.Info.CI)
	prj.Info.Homebrew = get("homebrew", prj.Info.Homebrew)
	prj.Info.Bundle = get("bundle", prj.Info.Bundle)
	prj.Info.SigningIdentity = get("signingIdentity", prj.Info.SigningIdentity)
	prj.Info.Slack = get("slack", prj.Info.Slack)
	prj.Info.JavaPackage = get("javaPackage", prj.Info.JavaPackage)
	prj.Info.GoVersion = get("goVersion", prj.Info.GoBinary)
	prj.Info.GoBinary = get("goBinary", prj.Info.GoBinary)
	prj.Info.ExtraFiles = util.StringSplitAndTrim(get("extraFiles", strings.Join(prj.Info.ExtraFiles, ", ")), ",")
	prj.Info.Deployments = util.StringSplitAndTrim(get("deployments", strings.Join(prj.Info.Deployments, ", ")), ",")

	var err error
	prj.Info.ModuleArgs, err = getModuleArgs(frm)
	if err != nil {
		return err
	}

	prj.Build = project.BuildFromMap(frm)
	if prj.Build.Empty() {
		prj.Build = nil
	}
	prj.Theme = theme.ApplyMap(frm)
	if prj.Theme.Equals(theme.ThemeDefault) {
		prj.Theme = nil
	}
	return nil
}

func getModuleArgs(frm util.ValueMap) (util.ValueMap, error) {
	ma := frm.GetStringOpt("moduleArgs")
	if ma == "" {
		return nil, nil
	}
	moduleArgs := util.ValueMap{}
	if err := util.FromJSON([]byte(ma), &moduleArgs); err != nil {
		return nil, errors.Wrap(err, "invalid module args JSON")
	}
	return moduleArgs, nil
}

func getProject(rc *fasthttp.RequestCtx, as *app.State) (*project.Project, error) {
	key, err := RCRequiredString(rc, "key", true)
	if err != nil {
		return nil, err
	}

	prj, err := as.Services.Projects.Get(key)
	if err != nil {
		return nil, err
	}
	return prj, nil
}
