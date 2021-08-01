package project

import (
	"fmt"
	"strings"

	"github.com/kyleu/projectforge/app/util"
)

type TemplateContext struct {
	Key     string `json:"key"`
	Type    string `json:"type"`
	Name    string `json:"name,omitempty"`
	Exec    string `json:"exec,omitempty"`
	Version string `json:"version"`
	Package string `json:"package,omitempty"`
	Args    string `json:"args,omitempty"`
	Port    int    `json:"port,omitempty"`

	Ignore     string `json:"ignore,omitempty"`
	IgnoreGrep string `json:"ignoreGrep,omitempty"`

	Org             string `json:"org,omitempty"`
	AuthorName      string `json:"authorName,omitempty"`
	AuthorEmail     string `json:"authorEmail,omitempty"`
	License         string `json:"license,omitempty"`
	Bundle          string `json:"bundle,omitempty"`
	SigningIdentity string `json:"signingIdentity,omitempty"`
	Homepage        string `json:"homepage,omitempty"`
	Sourcecode      string `json:"sourcecode,omitempty"`
	Summary         string `json:"summary,omitempty"`
	Description     string `json:"description,omitempty"`

	Build *Build `json:"build,omitempty"`
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

	ignore := strings.Join(util.StringArrayQuoted(p.Ignore), ", ")
	if ignore != "" {
		ignore = ", " + ignore
	}

	var ignoreGrep string
	for _, ig := range p.Ignore {
		ignoreGrep += fmt.Sprintf(" | grep -v \\\\./%s", ig)
	}

	ret := &TemplateContext{
		Key: p.Key, Type: p.Type, Name: p.Name, Exec: p.Exec, Version: p.Version, Package: p.Package, Args: p.Args, Port: p.Port,

		Ignore: ignore, IgnoreGrep: ignoreGrep,

		Org: i.Org, AuthorName: i.AuthorName, AuthorEmail: i.AuthorEmail, License: i.License, Bundle: i.Bundle, SigningIdentity: i.SigningIdentity,
		Homepage: i.Homepage, Sourcecode: i.Sourcecode, Summary: i.Summary, Description: i.Description,

		Build: b,
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
