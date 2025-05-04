package cproject

import (
	"fmt"
	"net/http"
	"net/url"
	"strconv"

	"projectforge.dev/projectforge/app"
	"projectforge.dev/projectforge/app/controller/cutil"
	"projectforge.dev/projectforge/app/lib/theme"
	"projectforge.dev/projectforge/app/project"
	"projectforge.dev/projectforge/app/util"
	"projectforge.dev/projectforge/views/layout"
	"projectforge.dev/projectforge/views/vpage"
)

const dblpipe = "||"

func projectFromForm(frm util.ValueMap, prj *project.Project) error {
	get := func(k string, def string) string {
		return util.OrDefault(frm.GetStringOpt(k), def)
	}
	prj.Name = get("name", prj.Name)
	prj.Icon = get("icon", prj.Icon)
	prj.Version = get("version", prj.Version)
	prj.Package = util.OrDefault(get("package", prj.Package), "github.com/org/"+prj.Key)
	prj.Args = get("args", prj.Args)
	prt, _ := strconv.ParseInt(get("port", fmt.Sprintf("%d", prj.Port)), 10, 32)
	prj.Port = int(prt)
	if prj.Port == 0 {
		prj.Port = 10000
	}
	prj.Modules = util.StringSplitAndTrim(get("modules", util.StringJoin(prj.Modules, dblpipe)), dblpipe)
	if len(prj.Modules) == 0 {
		prj.Modules = []string{"core"}
	}
	prj.Ignore = util.StringSplitAndTrim(get("ignore", util.StringJoin(prj.Ignore, ",")), ",")
	prj.Tags = util.StringSplitAndTrim(get("tags", util.StringJoin(prj.Tags, ",")), ",")
	prj.Path = get("path", prj.Path)

	if prj.Info == nil {
		prj.Info = &project.Info{}
	}
	prj.Info.Org = get("org", prj.Info.Org)
	if prj.Info.Org == "" {
		prj.Info.Org = prj.Key
		if prj.Info.Org == "" {
			prj.Info.Org = util.KeyUnknown
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
	prj.Info.Channels = util.StringSplitAndTrim(get("channels", util.StringJoin(prj.Info.Channels, ", ")), ",")
	prj.Info.JavaPackage = get("javaPackage", prj.Info.JavaPackage)
	prj.Info.GoVersion = get("goVersion", prj.Info.GoBinary)
	prj.Info.GoBinary = get("goBinary", prj.Info.GoBinary)
	prj.Info.ExtraFiles = util.StringSplitAndTrim(get("extraFiles", util.StringJoin(prj.Info.ExtraFiles, ", ")), ",")
	prj.Info.IgnoredFiles = util.StringSplitAndTrim(get("ignoredFiles", util.StringJoin(prj.Info.IgnoredFiles, ", ")), ",")
	prj.Info.Deployments = util.StringSplitAndTrim(get("deployments", util.StringJoin(prj.Info.Deployments, ", ")), ",")
	prj.Info.Acronyms = util.StringSplitAndTrim(get("acronyms", util.StringJoin(prj.Info.Acronyms, ", ")), ",")

	if prj.Package == "" {
		prj.Package = "github.com/" + prj.Info.Org + "/" + prj.Key
	}
	gh := "https://github.com/" + prj.Info.Org + "/" + prj.Key
	if prj.Info.Homepage == "" {
		prj.Info.Homepage = gh
	}
	if prj.Info.Sourcecode == "" {
		prj.Info.Sourcecode = gh
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

func getProject(r *http.Request, as *app.State) (*project.Project, error) {
	key, err := cutil.PathString(r, "key", true)
	if err != nil {
		return nil, err
	}

	prj, err := as.Services.Projects.Get(key)
	if err != nil {
		return nil, err
	}
	return prj, nil
}

func getProjectWithArgs(r *http.Request, as *app.State, logger util.Logger) (*project.Project, error) {
	prj, err := getProject(r, as)
	if err != nil {
		return nil, err
	}
	err = prj.ModuleArgExport(as.Services.Projects, logger)
	if err != nil {
		return nil, err
	}
	return prj, nil
}

func HandleLoad(cfg util.ValueMap, u *url.URL, title string) layout.Page {
	if cfg.GetStringOpt("hasloaded") == util.BoolTrue {
		return nil
	}
	cutil.URLAddQuery(u, "hasloaded", util.BoolTrue)
	page := &vpage.Load{URL: u.String(), Title: title}
	return page
}
