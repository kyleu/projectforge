package project

import (
	"fmt"
	"strings"

	"github.com/kyleu/projectforge/app/lib/theme"
	"github.com/kyleu/projectforge/app/util"
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

func (t *TemplateContext) Title() string {
	if t.Name != "" {
		return t.Name
	}
	return t.Key
}

func (t *TemplateContext) CleanKey() string {
	return clean(t.Key)
}

func (t *TemplateContext) KeyProper() string {
	return strings.ToUpper(t.Key[:1]) + t.Key[1:]
}

func (t *TemplateContext) NameCompressed() string {
	return strings.ReplaceAll(t.Name, " ", "")
}

func (t *TemplateContext) IgnoredSetting() string {
	if len(t.Ignore) == 0 {
		return ""
	}
	ret := make([]string, 0, len(t.Ignore))
	for _, i := range t.Ignore {
		ret = append(ret, "/"+i)
	}
	return " --skip-dirs \"" + strings.Join(ret, "|") + "\""
}

func (t *TemplateContext) IgnoredQuoted() string {
	if len(t.Ignore) == 0 {
		return ""
	}
	ret := make([]string, 0, len(t.Ignore))
	for _, i := range t.Ignore {
		ret = append(ret, fmt.Sprintf(", %q", i))
	}
	return strings.Join(ret, "")
}

func (t *TemplateContext) HasModule(m string) bool {
	return util.StringArrayContains(t.Modules, m)
}

func (t *TemplateContext) ModuleMarkdown() string {
	ret := make([]string, 0, len(t.Modules))
	for _, m := range t.Modules {
		ret = append(ret, fmt.Sprintf("- [%s](./doc/module/%s.md)", m, m))
	}
	return strings.Join(ret, "\n")
}

func (t *TemplateContext) PortIncremented(i int) int {
	return t.Port + i
}

func (t *TemplateContext) BuildAndroid() bool {
	ret := t.HasModule("android") && t.Build.Android
	return ret
}

func (t *TemplateContext) BuildIOS() bool {
	return t.HasModule("ios") && t.Build.IOS
}

func (t *TemplateContext) UsesLib() bool {
	return t.BuildIOS() || t.BuildAndroid() || t.Build.Desktop
}

func (t *TemplateContext) CIContent() string {
	if t.Info == nil {
		return ""
	}
	switch t.Info.CI {
	case "all":
		return "on: push"
	case "tags":
		return "on:\n  push:\n    tags"
	case "versions":
		return "on:\n  push:\n    tags:\n      - 'v*'"
	default:
		return "on:\n  push:\n    tags:\n      - 'DISABLED_v*'"
	}
}

func (t *TemplateContext) ExtraFilesContent() string {
	if t.Info == nil || len(t.Info.ExtraFiles) == 0 {
		return ""
	}
	ret := []string{"\n    extra_files:"}
	for _, ef := range t.Info.ExtraFiles {
		ret = append(ret, "      - "+ef)
	}
	return strings.Join(ret, "\n")
}

func (t *TemplateContext) ExtraFilesDocker() string {
	if t.Info == nil || len(t.Info.ExtraFiles) == 0 {
		return ""
	}
	ret := make([]string, 0, len(t.Info.ExtraFiles))
	for _, ef := range t.Info.ExtraFiles {
		ret = append(ret, "\nCOPY "+ef+" /")
	}
	return strings.Join(ret, "")
}

func (t *TemplateContext) HasSlack() bool {
	return t.Info.Slack != ""
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
