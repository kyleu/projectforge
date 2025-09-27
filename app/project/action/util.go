package action

import (
	"fmt"
	"slices"
	"strconv"
	"strings"
	"text/template"

	"github.com/pkg/errors"

	"projectforge.dev/projectforge/app/file"
	"projectforge.dev/projectforge/app/project"
	template2 "projectforge.dev/projectforge/app/project/template"
	"projectforge.dev/projectforge/app/util"
)

func projectFromCfg(proto *project.Project, cfg util.ValueMap) *project.Project {
	clean := func(s string) string {
		if strings.Contains(s, "\\") {
			s = util.StringSplitLastOnly(s, '\\', true)
		} else if strings.Contains(s, "/") {
			s = util.StringSplitLastOnly(s, '/', true)
		}
		return s
	}
	str := func(key string, def string) string {
		return util.OrDefault(cfg.GetStringOpt(key), def)
	}
	integer := func(key string, def int) int {
		s := str(key, "")
		i, err := strconv.ParseInt(s, 10, 32)
		if err != nil {
			return def
		}
		return int(i)
	}
	proto.Key = clean(proto.Key)
	proto.Name = clean(proto.Name)

	if proto.Package == "" {
		proto.Package = fmt.Sprintf("github.com/%s/%s", proto.Key, proto.Key)
	}

	port := integer("port", proto.Port)
	if port == 0 {
		port = 10000
	}

	mods, _ := cfg.GetStringArray("modules", true)
	if len(mods) == 1 {
		mods = util.StringSplitAndTrim(mods[0], "||")
	}
	if len(mods) == 0 {
		mods = []string{"core"}
	}
	slices.Sort(mods)

	return &project.Project{
		Key:        str("key", proto.Key),
		Name:       str("name", proto.Name),
		Icon:       str("icon", proto.Icon),
		Exec:       str("exec", proto.Icon),
		Version:    str("version", proto.Version),
		Package:    str("package", proto.Package),
		Args:       str("args", proto.Args),
		Port:       port,
		Modules:    mods,
		Ignore:     util.StringSplitAndTrim(str("ignore", util.StringJoin(proto.Ignore, ", ")), ","),
		Tags:       util.StringSplitAndTrim(str("tags", util.StringJoin(proto.Tags, ", ")), ","),
		Info:       infoFromCfg(proto, cfg),
		Theme:      proto.Theme,
		Build:      proto.Build,
		ExportArgs: proto.ExportArgs,
		Config:     proto.Config,
		Path:       proto.Path,
		Parent:     proto.Parent,
	}
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

	docs := i.Docs
	if x := cfg.GetStringOpt("docs"); x != "" {
		_ = util.FromJSON([]byte(x), &docs)
	}

	additionalPorts := i.AdditionalPorts
	if x := cfg.GetStringOpt("additionalPorts"); x != "" {
		_ = util.FromJSON([]byte(x), &additionalPorts)
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
		NotarizeEmail:   str("notarizeEmail", i.NotarizeEmail),
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
