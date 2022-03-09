package project

import (
	"fmt"

	"projectforge.dev/projectforge/app/lib/theme"
)

type TemplateContext struct {
	Key         string         `json:"key"`
	Name        string         `json:"name,omitempty"`
	Exec        string         `json:"exec,omitempty"`
	Version     string         `json:"version"`
	Package     string         `json:"package,omitempty"`
	Args        string         `json:"args,omitempty"`
	Port        int            `json:"port,omitempty"`
	PortOffsets map[string]int `json:"portOffsets,omitempty"`

	Modules []string     `json:"modules,omitempty"`
	Info    *Info        `json:"info,omitempty"`
	Build   *Build       `json:"build,omitempty"`
	Theme   *theme.Theme `json:"theme,omitempty"`

	Ignore     []string `json:"ignore,omitempty"`
	IgnoreGrep string   `json:"ignoreGrep,omitempty"`
}

func (p *Project) ToTemplateContext(portOffsets map[string]int) *TemplateContext {
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
		ignoreGrep += fmt.Sprintf(" | grep -v \\\\./%s", ig)
	}

	ret := &TemplateContext{
		Key: p.Key, Name: p.Name, Exec: p.Executable(), Version: p.Version,
		Package: p.Package, Args: p.Args, Port: p.Port, PortOffsets: portOffsets,
		Modules: p.Modules, Info: i, Build: b, Theme: t, Ignore: p.Ignore, IgnoreGrep: ignoreGrep,
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
