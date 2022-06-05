package action

import (
	"strconv"
	"strings"
	"text/template"

	"github.com/pkg/errors"

	"projectforge.dev/projectforge/app/file"
	"projectforge.dev/projectforge/app/project"
	"projectforge.dev/projectforge/app/util"
)

func projectFromCfg(proto *project.Project, cfg util.ValueMap) *project.Project {
	str := func(key string, def string) string {
		ret := cfg.GetStringOpt(key)
		if ret == "" {
			return def
		}
		return ret
	}
	integer := func(key string, def int) int {
		s := str(key, "")
		i, err := strconv.Atoi(s)
		if err != nil {
			return def
		}
		return i
	}

	mods, _ := cfg.GetStringArray("modules", true)
	if len(mods) == 0 {
		mods = []string{"core"}
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

	return &project.Project{
		Key:     str("key", proto.Key),
		Version: str("version", proto.Version),
		Name:    str("name", proto.Name),
		Package: str("package", proto.Package),
		Args:    str("args", proto.Args),
		Port:    integer("port", proto.Port),
		Modules: mods,
		Ignore:  util.StringSplitAndTrim(str("ignore", strings.Join(proto.Ignore, ", ")), ","),
		Tags:    util.StringSplitAndTrim(str("tags", strings.Join(proto.Tags, ", ")), ","),
		Info: &project.Info{
			Org:             str("org", i.Org),
			AuthorID:        str("author_id", i.AuthorID),
			AuthorName:      str("author_name", i.AuthorName),
			AuthorEmail:     str("author_email", i.AuthorEmail),
			License:         str("license", i.License),
			Homepage:        str("homepage", i.Homepage),
			Sourcecode:      str("sourcecode", i.Sourcecode),
			Summary:         str("summary", i.Summary),
			Description:     str("description", i.Description),
			CI:              str("ci", i.CI),
			Homebrew:        str("homebrew", i.Homebrew),
			Bundle:          str("bundle", i.Bundle),
			SigningIdentity: str("signingIdentity", i.SigningIdentity),
			Slack:           str("slack", i.Slack),
			JavaPackage:     str("javaPackage", i.JavaPackage),
			GoVersion:       str("goVersion", i.GoBinary),
			GoBinary:        str("goBinary", i.GoBinary),
			ConfigVars:      cfgVars,
			ExtraFiles:      util.StringSplitAndTrim(str("extraFiles", strings.Join(i.ExtraFiles, ", ")), ","),
			Deployments:     util.StringSplitAndTrim(str("deployments", strings.Join(i.Deployments, ", ")), ","),
			ModuleDefs:      md,
		},
		Path: proto.Path,
	}
}

func runTemplate(path string, content string, ctx *project.TemplateContext) (string, error) {
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

func runTemplateFile(f *file.File, ctx *project.TemplateContext) (string, error) {
	return runTemplate(f.FullPath(), f.Content, ctx)
}
