package template

import (
	"fmt"
	"strings"

	"github.com/samber/lo"

	"projectforge.dev/projectforge/app/lib/metamodel/model"
	"projectforge.dev/projectforge/app/lib/theme"
	"projectforge.dev/projectforge/app/project"
	"projectforge.dev/projectforge/app/util"
)

type Context struct {
	Key         string            `json:"key"`
	Name        string            `json:"name,omitempty"`
	Exec        string            `json:"exec,omitempty"`
	Version     string            `json:"version"`
	Package     string            `json:"package,omitempty"`
	Args        string            `json:"args,omitempty"`
	Port        int               `json:"port,omitempty"`
	ConfigVars  util.KeyTypeDescs `json:"configVars,omitempty"`
	PortOffsets map[string]int    `json:"portOffsets,omitempty"`

	Ignore     []string `json:"ignore,omitempty"`
	IgnoreGrep string   `json:"ignoreGrep,omitempty"`
	Modules    []string `json:"modules,omitempty"`
	Tags       []string `json:"tags,omitempty"`

	ExportArgs     *model.Args    `json:"exportArgs,omitempty"`
	Config         util.ValueMap  `json:"config,omitempty"`
	Info           *project.Info  `json:"info,omitempty"`
	Build          *project.Build `json:"build,omitempty"`
	Theme          *theme.Theme   `json:"theme,omitempty"`
	DatabaseEngine string         `json:"databaseEngine,omitempty"`
	Linebreak      string         `json:"-"`
}

func ToTemplateContext(p *project.Project, configVars util.KeyTypeDescs, portOffsets map[string]int, linebreak string) *Context {
	i := p.Info
	if i == nil {
		i = &project.Info{}
	}
	b := p.Build
	if b == nil {
		b = &project.Build{}
	}
	t := p.Theme
	if t == nil {
		t = theme.Default
	}

	var ignoreGrep string
	lo.ForEach(p.Ignore, func(ig string, _ int) {
		if strings.HasPrefix(ig, "^") {
			ig = "\\\\./" + strings.TrimPrefix(ig, "^")
		}
		ignoreGrep += fmt.Sprintf(" | grep -v %s", ig)
	})

	cv := append(util.KeyTypeDescs{}, configVars...)
	if p.Info != nil && len(p.Info.ConfigVars) > 0 {
		cv = append(cv, p.Info.ConfigVars...)
	}

	ret := &Context{
		Key: p.Key, Name: p.Name, Exec: p.Executable(), Version: p.Version,
		Package: p.Package, Args: p.Args, Port: p.Port, ConfigVars: cv, PortOffsets: portOffsets,
		Ignore: p.Ignore, IgnoreGrep: ignoreGrep, Modules: p.Modules, Tags: p.Tags,
		ExportArgs: p.ExportArgs, Config: p.Config, Info: i, Build: b, Theme: t, DatabaseEngine: p.DatabaseEngineDefault(), Linebreak: linebreak,
	}

	if ret.Name == "" {
		ret.Name = ret.Key
	}
	if ret.Args == "" {
		if p.HasModule("marketing") {
			ret.Args = " -v --addr=0.0.0.0 all"
		} else {
			ret.Args = " -v --addr=0.0.0.0 server"
		}
	}

	return ret
}
