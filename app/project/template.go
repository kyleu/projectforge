package project

import (
	"fmt"
	"strings"

	"projectforge.dev/projectforge/app/lib/theme"
	"projectforge.dev/projectforge/app/util"
)

type TemplateContext struct {
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

	ModuleArgs util.ValueMap `json:"moduleArgs,omitempty"`
	Info       *Info         `json:"info,omitempty"`
	Build      *Build        `json:"build,omitempty"`
	Theme      *theme.Theme  `json:"theme,omitempty"`
}

func (p *Project) ToTemplateContext(configVars util.KeyTypeDescs, portOffsets map[string]int) *TemplateContext {
	i := p.Info
	if i == nil {
		i = &Info{}
	}
	b := p.Build
	if b == nil {
		b = &Build{}
	}
	t := p.Theme
	if t == nil {
		t = theme.ThemeDefault
	}

	var ignoreGrep string
	for _, ig := range p.Ignore {
		if strings.HasPrefix(ig, "^") {
			ig = "\\\\./" + strings.TrimPrefix(ig, "^")
		}
		ignoreGrep += fmt.Sprintf(" | grep -v %s", ig)
	}

	cv := append(util.KeyTypeDescs{}, configVars...)
	if p.Info != nil && len(p.Info.ConfigVars) > 0 {
		cv = append(cv, p.Info.ConfigVars...)
	}

	ret := &TemplateContext{
		Key: p.Key, Name: p.Name, Exec: p.Executable(), Version: p.Version,
		Package: p.Package, Args: p.Args, Port: p.Port, ConfigVars: cv, PortOffsets: portOffsets,
		Ignore: p.Ignore, IgnoreGrep: ignoreGrep, Modules: p.Modules, Tags: p.Tags,
		ModuleArgs: p.ModuleArgs, Info: i, Build: b, Theme: t,
	}

	if ret.Name == "" {
		ret.Name = ret.Key
	}
	if ret.Exec == "" {
		ret.Exec = ret.Key
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
