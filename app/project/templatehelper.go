package project

import (
	"fmt"
	"strings"

	"github.com/kyleu/projectforge/app/util"
)

const goStdBin = "go"

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

func (t *TemplateContext) HasSlack() bool {
	return t.Info.Slack != ""
}

func (t *TemplateContext) GoVersionSafe() string {
	if t.Info.GoVersion == "" {
		return defaultGoVersion
	}
	return t.Info.GoVersion + "\n          stable: false"
}

func (t *TemplateContext) GoBinarySafe() string {
	if t.Info.GoBinary == "" {
		return goStdBin
	}
	return t.Info.GoBinary
}
