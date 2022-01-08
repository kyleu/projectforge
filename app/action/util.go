package action

import (
	"strconv"

	"github.com/kyleu/projectforge/app/project"
	"github.com/kyleu/projectforge/app/util"
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

	var moduleArgs util.ValueMap
	if ma, err := cfg.GetMap("moduleArgs"); ma != nil && err == nil {
		moduleArgs = ma
	}

	return &project.Project{
		Key:     str("key", proto.Key),
		Version: str("version", proto.Version),
		Name:    str("name", proto.Name),
		Package: str("package", proto.Package),
		Args:    str("args", proto.Args),
		Port:    integer("port", proto.Port),
		Modules: mods,
		Ignore:  proto.Ignore,
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
			ModuleArgs:      moduleArgs,
		},
		Path: proto.Path,
	}
}
