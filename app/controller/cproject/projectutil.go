package cproject

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/valyala/fasthttp"

	"projectforge.dev/projectforge/app"
	"projectforge.dev/projectforge/app/controller/cutil"
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
	prj.Icon = get("icon", prj.Icon)
	prj.Version = get("version", prj.Version)
	prj.Package = get("package", prj.Package)
	if prj.Package == "" {
		prj.Package = "github.com/org/" + prj.Key
	}
	prj.Args = get("args", prj.Args)
	prt, _ := strconv.ParseInt(get("port", fmt.Sprintf("%d", prj.Port)), 10, 32)
	prj.Port = int(prt)
	if prj.Port == 0 {
		prj.Port = 10000
	}
	prj.Modules = util.StringSplitAndTrim(get("modules", strings.Join(prj.Modules, "||")), "||")
	if len(prj.Modules) == 0 {
		prj.Modules = []string{"core"}
	}
	prj.Ignore = util.StringSplitAndTrim(get("ignore", strings.Join(prj.Ignore, ",")), ",")
	prj.Tags = util.StringSplitAndTrim(get("tags", strings.Join(prj.Tags, ",")), ",")
	prj.Path = get("path", prj.Path)

	if prj.Info == nil {
		prj.Info = &project.Info{}
	}
	prj.Info.Org = get("org", prj.Info.Org)
	if prj.Info.Org == "" {
		prj.Info.Org = prj.Key
		if prj.Info.Org == "" {
			prj.Info.Org = "unknown"
		}
	}

	prj.Info.AuthorID = get("authorID", prj.Info.AuthorID)
	prj.Info.AuthorName = get("authorName", prj.Info.AuthorName)
	prj.Info.AuthorEmail = get("authorEmail", prj.Info.AuthorEmail)
	prj.Info.Team = get("team", prj.Info.Team)
	prj.Info.License = get("license", prj.Info.License)
	prj.Info.Homepage = get("homepage", prj.Info.Homepage)
	prj.Info.Sourcecode = get("sourcecode", prj.Info.Sourcecode)
	prj.Info.Summary = get("summary", prj.Info.Summary)
	prj.Info.Description = get("description", prj.Info.Description)
	prj.Info.CI = get("ci", prj.Info.CI)
	prj.Info.Homebrew = get("homebrew", prj.Info.Homebrew)
	prj.Info.Bundle = get("bundle", prj.Info.Bundle)
	prj.Info.SigningIdentity = get("signingIdentity", prj.Info.SigningIdentity)
	prj.Info.NotarizeEmail = get("notarizeEmail", prj.Info.NotarizeEmail)
	prj.Info.DatabaseEngine = get("databaseEngine", prj.Info.DatabaseEngine)
	prj.Info.Slack = get("slack", prj.Info.Slack)
	prj.Info.Channels = util.StringSplitAndTrim(get("channels", strings.Join(prj.Info.Channels, ", ")), ",")
	prj.Info.JavaPackage = get("javaPackage", prj.Info.JavaPackage)
	prj.Info.GoVersion = get("goVersion", prj.Info.GoBinary)
	prj.Info.GoBinary = get("goBinary", prj.Info.GoBinary)
	prj.Info.ExtraFiles = util.StringSplitAndTrim(get("extraFiles", strings.Join(prj.Info.ExtraFiles, ", ")), ",")
	prj.Info.Deployments = util.StringSplitAndTrim(get("deployments", strings.Join(prj.Info.Deployments, ", ")), ",")
	prj.Info.Acronyms = util.StringSplitAndTrim(get("acronyms", strings.Join(prj.Info.Acronyms, ", ")), ",")

	if prj.Package == "" {
		prj.Package = "github.com/" + prj.Info.Org + "/" + prj.Key
	}
	if prj.Info.Homepage == "" {
		prj.Info.Homepage = "https://github.com/" + prj.Info.Org + "/" + prj.Key
	}
	if prj.Info.Sourcecode == "" {
		prj.Info.Sourcecode = "https://github.com/" + prj.Info.Org + "/" + prj.Key
	}

	prj.Build = project.BuildFromMap(frm)
	if prj.Build.Empty() {
		prj.Build = nil
	}
	prj.Theme = theme.ApplyMap(frm)
	if prj.Theme.Equals(theme.Default) {
		prj.Theme = nil
	}
	return nil
}

func getProject(rc *fasthttp.RequestCtx, as *app.State) (*project.Project, error) {
	key, err := cutil.RCRequiredString(rc, "key", true)
	if err != nil {
		return nil, err
	}

	prj, err := as.Services.Projects.Get(key)
	if err != nil {
		return nil, err
	}
	return prj, nil
}
