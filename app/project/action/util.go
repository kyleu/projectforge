package action

import (
	"strings"
	"text/template"

	"github.com/pkg/errors"

	"projectforge.dev/projectforge/app/file"
	"projectforge.dev/projectforge/app/lib/theme"
	"projectforge.dev/projectforge/app/project"
	template2 "projectforge.dev/projectforge/app/project/template"
	"projectforge.dev/projectforge/app/util"
)

func ProjectFromMap(prj *project.Project, m util.ValueMap, parseKey bool) error {
	clean := func(s string) string {
		if strings.Contains(s, "\\") {
			s = util.StringSplitLastOnly(s, '\\', true)
		} else if strings.Contains(s, "/") {
			s = util.StringSplitLastOnly(s, '/', true)
		}
		return s
	}
	get := func(k string, def string) util.Str {
		return m.GetRichStringOpt(k).OrDefault(def)
	}
	getSplit := func(k string, def []string) util.Strings {
		delim := "||"
		return get(k, util.StringJoin(def, delim)).SplitAndTrim(delim)
	}
	getInt := func(key string, def int) int {
		i, ok := get(key, "").ParseInt()
		return util.Choose(ok, i, def)
	}

	if parseKey {
		prj.Key = clean(string(get("key", prj.Name).OrDefault(prj.Key)))
	}

	prj.Name = clean(get("name", prj.Name).String())
	prj.Icon = get("icon", prj.Icon).String()
	prj.Exec = get("exec", prj.Exec).String()
	prj.Version = get("version", prj.Version).String()
	prj.Package = get("package", prj.Package).OrDefault("github.com/" + prj.Key + "/" + prj.Key).String()
	prj.Args = get("args", prj.Args).String()
	prj.Port = getInt("port", util.Choose(prj.Port == 0, 10000, prj.Port))
	prj.Modules = util.ArraySorted(getSplit("modules", prj.Modules).Strings())
	if len(prj.Modules) == 0 {
		prj.Modules = []string{"core"}
	}
	prj.Ignore = getSplit("ignore", prj.Ignore).Strings()
	prj.Tags = getSplit("tags", prj.Tags).Strings()

	prj.Info = infoFromCfg(prj, m)

	if _, ok := m["light-background"]; ok {
		prj.Theme = theme.ApplyMap(prj.Theme, m)
	}
	if prj.Theme.Equals(theme.Default) {
		prj.Theme = nil
	}

	prj.Build = project.BuildFromMap(m)
	if prj.Build.Empty() {
		prj.Build = nil
	}

	prj.Files = getSplit("files", prj.Files).Strings()

	return nil
}

func infoFromCfg(proto *project.Project, cfg util.ValueMap) *project.Info {
	str := func(key string, def string) string {
		return util.OrDefault(cfg.GetStringOpt(key), def)
	}

	i := proto.Info
	if i == nil {
		i = &project.Info{License: "Proprietary"}
	}

	md := i.ModuleDefs
	if x := cfg.GetStringOpt("moduleDefs"); x != "" {
		_ = util.FromJSON([]byte(x), &md)
	}

	cfgVars := i.ConfigVars
	if x := cfg.GetStringOpt("configVars"); x != "" {
		_ = util.FromJSON([]byte(x), &cfgVars)
	}

	deps := i.Dependencies
	if x := cfg.GetStringOpt("dependencies"); x != "" {
		_ = util.FromJSON([]byte(x), &deps)
	}

	docs := i.Docs
	if x := cfg.GetStringOpt("docs"); x != "" {
		_ = util.FromJSON([]byte(x), &docs)
	}

	additionalPorts := i.AdditionalPorts
	if x := cfg.GetStringOpt("additionalPorts"); x != "" {
		_ = util.FromJSON([]byte(x), &additionalPorts)
		if len(additionalPorts) == 0 {
			additionalPorts = nil
		}
	}

	return &project.Info{
		Org:             str("org", i.Org),
		AuthorID:        str("author_id", i.AuthorID),
		AuthorName:      str("author_name", i.AuthorName),
		AuthorEmail:     str("author_email", i.AuthorEmail),
		Team:            str("team", i.Team),
		License:         str("license", i.License),
		Homepage:        str("homepage", i.Homepage),
		Sourcecode:      str("sourcecode", i.Sourcecode),
		Summary:         str("summary", i.Summary),
		Description:     str("description", i.Description),
		CI:              str("ci", i.CI),
		Homebrew:        str("homebrew", i.Homebrew),
		Bundle:          str("bundle", i.Bundle),
		SigningIdentity: str("signingIdentity", i.SigningIdentity),
		DatabaseEngine:  str("databaseEngine", i.DatabaseEngine),
		Slack:           str("slack", i.Slack),
		Channels:        util.StringSplitAndTrim(str("channels", util.StringJoin(i.Channels, ", ")), ","),
		JavaPackage:     str("javaPackage", i.JavaPackage),
		GoVersion:       str("goVersion", i.GoVersion),
		GoBinary:        str("goBinary", i.GoBinary),
		ConfigVars:      cfgVars,
		AdditionalPorts: additionalPorts,
		ExtraFiles:      util.StringSplitAndTrim(str("extraFiles", util.StringJoin(i.ExtraFiles, ", ")), ","),
		IgnoredFiles:    util.StringSplitAndTrim(str("ignoredFiles", util.StringJoin(i.IgnoredFiles, ", ")), ","),
		Deployments:     util.StringSplitAndTrim(str("deployments", util.StringJoin(i.Deployments, ", ")), ","),
		EnvVars:         util.StringSplitAndTrim(str("envvars", util.StringJoin(i.EnvVars, ", ")), ","),
		DockerPackages:  util.StringSplitAndTrim(str("dockerPackages", util.StringJoin(i.DockerPackages, ", ")), ","),
		Dependencies:    deps,
		Docs:            docs,
		Acronyms:        util.StringSplitAndTrim(str("acronyms", util.StringJoin(i.Acronyms, ", ")), ","),
		ModuleDefs:      md,
	}
}

func runTemplate(path string, content string, ctx *template2.Context) (string, error) {
	t, err := template.New(path).Delims(delimStart, delimEnd).Parse(content)
	if err != nil {
		return "", errors.Wrapf(err, "unable to create template for [%s]", path)
	}

	res := &strings.Builder{}
	err = t.Execute(res, ctx)
	if err != nil {
		return "", errors.Wrapf(err, "unable to execute template for [%s]", path)
	}
	return res.String(), nil
}

func runTemplateFile(f *file.File, ctx *template2.Context) (string, error) {
	return runTemplate(f.FullPath(), f.Content, ctx)
}
