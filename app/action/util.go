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
	i := proto.Info
	if i == nil {
		i = &project.Info{License: "Proprietary"}
	}
	return &project.Project{
		Key:     str("key", proto.Key),
		Version: proto.Version,
		Name:    str("name", proto.Name),
		Package: str("package", proto.Package),
		Args:    str("args", proto.Args),
		Port:    integer("port", proto.Port),
		Modules: []string{"core"},
		Ignore:  proto.Ignore,
		Info: &project.Info{
			Org:             str("org", i.Org),
			AuthorName:      str("author_name", i.AuthorName),
			AuthorEmail:     str("author_email", i.AuthorEmail),
			License:         str("license", i.License),
			Homepage:        str("homepage", i.Homepage),
			Sourcecode:      str("sourcecode", i.Sourcecode),
			Summary:         str("summary", i.Summary),
			Description:     str("description", i.Description),
			CI:              str("ci", i.Description),
			Homebrew:        str("homebrew", i.Description),
			Bundle:          str("bundle", i.Description),
			SigningIdentity: str("signingIdentity", i.Description),
			Slack:           str("slack", i.Description),
			JavaPackage:     str("javaPackage", i.Description),
		},
		Path: proto.Path,
	}
}
