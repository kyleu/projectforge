package project

import (
	"sort"

	"github.com/kyleu/projectforge/app/theme"
	"github.com/kyleu/projectforge/app/util"
)

var DefaultIcon = "code"

type Project struct {
	Key      string       `json:"key"`
	Name     string       `json:"name,omitempty"`
	Icon     string       `json:"icon,omitempty"`
	Exec     string       `json:"exec,omitempty"`
	Version  string       `json:"version"`
	Package  string       `json:"package,omitempty"`
	Args     string       `json:"args,omitempty"`
	Port     int          `json:"port,omitempty"`
	Modules  []string     `json:"modules"`
	Ignore   []string     `json:"ignore,omitempty"`
	Children []string     `json:"children,omitempty"`
	Info     *Info        `json:"info,omitempty"`
	Theme    *theme.Theme `json:"theme,omitempty"`
	Build    *Build       `json:"build,omitempty"`

	Path   string `json:"-"`
	Parent string `json:"-"`
}

func (p *Project) Title() string {
	if p.Name == "" {
		return p.Key
	}
	return p.Name
}

func (p *Project) Executable() string {
	if p.Exec == "" {
		return p.Key
	}
	return p.Exec
}

func (p *Project) SafeIcon() string {
	_, ok := util.SVGLibrary[p.Icon]
	if !ok {
		return DefaultIcon
	}
	return p.Icon
}

func NewProject(key string, path string) *Project {
	return &Project{Key: key, Version: "0.0.0", Path: path}
}

type Projects []*Project

func (p Projects) AllModules() []string {
	var ret []string
	for _, prj := range p {
		for _, mod := range prj.Modules {
			hit := false
			for _, x := range ret {
				if x == mod {
					hit = true
					break
				}
			}
			if !hit {
				ret = append(ret, mod)
			}
		}
	}
	sort.Strings(ret)
	return ret
}
