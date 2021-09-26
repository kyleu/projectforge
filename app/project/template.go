package project

import (
	"fmt"
	"strings"

	"github.com/kyleu/projectforge/app/theme"
	"github.com/kyleu/projectforge/app/util"
)

type TemplateContext struct {
	Key     string `json:"key"`
	Name    string `json:"name,omitempty"`
	Exec    string `json:"exec,omitempty"`
	Version string `json:"version"`
	Package string `json:"package,omitempty"`
	Args    string `json:"args,omitempty"`
	Port    int    `json:"port,omitempty"`

	Modules []string     `json:"modules,omitempty"`
	Info    *Info        `json:"info,omitempty"`
	Build   *Build       `json:"build,omitempty"`
	Theme   *theme.Theme `json:"theme,omitempty"`

	Ignore     string `json:"ignore,omitempty"`
	IgnoreGrep string `json:"ignoreGrep,omitempty"`
}

func (t *TemplateContext) HasModule(m string) bool {
	return util.StringArrayContains(t.Modules, m)
}

func (t *TemplateContext) BuildAndroid() bool {
	ret := t.HasModule("android") && t.Build.Android
	return ret
}

func (t *TemplateContext) BuildIOS() bool {
	return t.HasModule("ios") && t.Build.IOS
}

func (p *Project) ToTemplateContext() *TemplateContext {
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

	ignore := strings.Join(util.StringArrayQuoted(p.Ignore), ", ")
	if ignore != "" {
		ignore = ", " + ignore
	}

	var ignoreGrep string
	for _, ig := range p.Ignore {
		ignoreGrep += fmt.Sprintf(" | grep -v \\\\./%s", ig)
	}

	ret := &TemplateContext{
		Key: p.Key, Name: p.Name, Exec: p.Exec, Version: p.Version, Package: p.Package, Args: p.Args, Port: p.Port,
		Modules: p.Modules, Info: i, Build: b, Theme: t, Ignore: ignore, IgnoreGrep: ignoreGrep,
	}

	if ret.Name == "" {
		ret.Name = ret.Key
	}
	if ret.Exec == "" {
		ret.Exec = ret.Key
	}
	if ret.Args == "" {
		if util.StringArrayContains(p.Modules, "marketing") {
			ret.Args = " -v --addr=0.0.0.0 all"
		} else {
			ret.Args = " -v --addr=0.0.0.0 server"
		}
	}

	return ret
}
