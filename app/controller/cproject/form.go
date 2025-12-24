package cproject

import (
	"fmt"

	"projectforge.dev/projectforge/app/lib/theme"
	"projectforge.dev/projectforge/app/project"
	"projectforge.dev/projectforge/app/util"
)

func projectFromForm(frm util.ValueMap, prj *project.Project) error {
	get := func(k string, def string) util.Str {
		return frm.GetRichStringOpt(k).OrDefault(def)
	}
	getSplit := func(k string, def []string, delimOpt ...string) util.Strings {
		delim := ","
		if len(delimOpt) > 0 {
			delim = delimOpt[0]
		}
		return get(k, util.StringJoin(def, delim)).SplitAndTrim(delim)
	}
	prj.Name = get("name", prj.Name).String()
	prj.Icon = get("icon", prj.Icon).String()
	prj.Version = get("version", prj.Version).String()
	prj.Package = get("package", prj.Package).OrDefault("github.com/org/" + prj.Key).String()
	prj.Args = get("args", prj.Args).String()
	prt, _ := get("port", fmt.Sprintf("%d", prj.Port)).ParseInt()
	prj.Port = prt
	if prj.Port == 0 {
		prj.Port = 10000
	}
	prj.Modules = getSplit("modules", prj.Modules, dblpipe).Strings()
	if len(prj.Modules) == 0 {
		prj.Modules = []string{"core"}
	}
	prj.Ignore = getSplit("ignore", prj.Ignore).Strings()
	prj.Tags = getSplit("tags", prj.Tags).Strings()
	prj.Path = get("path", prj.Path).String()

	if prj.Info == nil {
		prj.Info = &project.Info{}
	}
	prj.Info.Org = get("org", prj.Info.Org).String()
	if prj.Info.Org == "" {
		prj.Info.Org = prj.Key
		if prj.Info.Org == "" {
			prj.Info.Org = util.KeyUnknown
		}
	}

	prj.Info.AuthorID = get("authorID", prj.Info.AuthorID).String()
	prj.Info.AuthorName = get("authorName", prj.Info.AuthorName).String()
	prj.Info.AuthorEmail = get("authorEmail", prj.Info.AuthorEmail).String()
	prj.Info.Team = get("team", prj.Info.Team).String()
	prj.Info.License = get("license", prj.Info.License).String()
	prj.Info.Homepage = get("homepage", prj.Info.Homepage).String()
	prj.Info.Sourcecode = get("sourcecode", prj.Info.Sourcecode).String()
	prj.Info.Summary = get("summary", prj.Info.Summary).String()
	prj.Info.Description = get("description", prj.Info.Description).String()
	prj.Info.CI = get("ci", prj.Info.CI).String()
	prj.Info.Homebrew = get("homebrew", prj.Info.Homebrew).String()
	prj.Info.Bundle = get("bundle", prj.Info.Bundle).String()
	prj.Info.SigningIdentity = get("signingIdentity", prj.Info.SigningIdentity).String()
	prj.Info.DatabaseEngine = get("databaseEngine", prj.Info.DatabaseEngine).String()
	prj.Info.Slack = get("slack", prj.Info.Slack).String()
	prj.Info.Channels = getSplit("channels", prj.Info.Channels).Strings()
	prj.Info.JavaPackage = get("javaPackage", prj.Info.JavaPackage).String()
	prj.Info.GoVersion = get("goVersion", prj.Info.GoVersion).String()
	prj.Info.GoBinary = get("goBinary", prj.Info.GoBinary).String()
	prj.Info.ExtraFiles = getSplit("extraFiles", prj.Info.ExtraFiles).Strings()
	prj.Info.IgnoredFiles = getSplit("ignoredFiles", prj.Info.IgnoredFiles).Strings()
	prj.Info.Deployments = getSplit("deployments", prj.Info.Deployments).Strings()
	prj.Info.Acronyms = getSplit("acronyms", prj.Info.Acronyms).Strings()

	cv := get("configVars", util.ToJSON(prj.Info.ConfigVars)).String()
	if err := util.FromJSON([]byte(cv), &prj.Info.ConfigVars); err != nil {
		return err
	}
	if len(prj.Info.ConfigVars) == 0 {
		prj.Info.ConfigVars = nil
	}
	ap := get("additionalPorts", util.ToJSON(prj.Info.AdditionalPorts))
	if err := util.FromJSON([]byte(ap), &prj.Info.AdditionalPorts); err != nil {
		return err
	}
	if len(prj.Info.AdditionalPorts) == 0 {
		prj.Info.AdditionalPorts = nil
	}

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
	prj.Theme = theme.ApplyMap(prj.Theme, frm)
	if prj.Theme.Equals(theme.Default) {
		prj.Theme = nil
	}
	return nil
}
