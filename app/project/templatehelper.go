package project

import (
	"fmt"
	"strings"

	"golang.org/x/exp/slices"
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

func (t *TemplateContext) SourceTrimmed() string {
	return strings.TrimPrefix(strings.TrimPrefix(t.Info.Sourcecode, "http://"), "https://")
}

func (t *TemplateContext) HasModule(m string) bool {
	return slices.Contains(t.Modules, m)
}

func (t *TemplateContext) HasDatabaseModule() bool {
	return t.HasModule("migration") || t.HasModule("readonlydb")
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

func (t *TemplateContext) BuildMobile() bool {
	return t.BuildIOS() || t.BuildAndroid()
}

func (t *TemplateContext) UsesLib() bool {
	return t.BuildMobile() || t.Build.Desktop
}

func (t *TemplateContext) HasSlack() bool {
	return t.Info.Slack != ""
}

func (t *TemplateContext) DatabaseUIOpts() (bool, bool) {
	cfg, _ := t.Info.ModuleArgs.GetMap("databaseui", true)
	if len(cfg) == 0 {
		return true, true
	}
	sqleditor, _ := cfg.GetBool("sqleditor", true)
	readonly, _ := cfg.GetBool("readonly", true)
	return sqleditor, readonly
}

func (t *TemplateContext) DatabaseUISQLEditor() bool {
	sqleditor, _ := t.DatabaseUIOpts()
	return sqleditor
}

func (t *TemplateContext) DatabaseUIReadOnly() bool {
	_, readonly := t.DatabaseUIOpts()
	return readonly
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
