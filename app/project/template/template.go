package template

import (
	"fmt"
	"strings"

	"github.com/samber/lo"

	"projectforge.dev/projectforge/app/lib/metamodel"
	"projectforge.dev/projectforge/app/lib/theme"
	"projectforge.dev/projectforge/app/project"
	"projectforge.dev/projectforge/app/util"
)

type Context struct {
	Key         string            `json:"key"`
	Name        string            `json:"name,omitzero"`
	Exec        string            `json:"exec,omitzero"`
	Version     string            `json:"version"`
	Package     string            `json:"package,omitzero"`
	Args        []string          `json:"args,omitempty"`
	Port        int               `json:"port,omitzero"`
	ConfigVars  util.KeyTypeDescs `json:"configVars,omitzero"`
	PortOffsets map[string]int    `json:"portOffsets,omitzero"`

	Ignore     []string `json:"ignore,omitempty"`
	IgnoreGrep string   `json:"ignoreGrep,omitzero"`
	Modules    []string `json:"modules,omitempty"`
	Tags       []string `json:"tags,omitempty"`

	ExportArgs     *metamodel.Args `json:"exportArgs,omitzero"`
	Config         util.ValueMap   `json:"config,omitzero"`
	Info           *project.Info   `json:"info,omitzero"`
	Build          *project.Build  `json:"build,omitzero"`
	Theme          *theme.Theme    `json:"theme,omitzero"`
	DatabaseEngine string          `json:"databaseEngine,omitzero"`
	Linebreak      string          `json:"-"`
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

	args := util.StringSplitAndTrim(p.Args, " ")

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
		Package: p.Package, Args: args, Port: p.Port, ConfigVars: cv, PortOffsets: portOffsets,
		Ignore: p.Ignore, IgnoreGrep: ignoreGrep, Modules: p.Modules, Tags: p.Tags,
		ExportArgs: p.ExportArgs, Config: p.Config, Info: i, Build: b, Theme: t, DatabaseEngine: p.DatabaseEngineDefault(), Linebreak: linebreak,
	}

	if ret.Name == "" {
		ret.Name = ret.Key
	}
	if len(ret.Args) == 0 {
		if p.HasModule("marketing") {
			ret.Args = []string{"-v", "--addr=0.0.0.0", "all"}
		} else {
			ret.Args = []string{"-v", "--addr=0.0.0.0", "server"}
		}
	}

	return ret
}
