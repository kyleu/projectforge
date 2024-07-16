package project

import (
	"fmt"
	"strings"

	"github.com/samber/lo"

	"projectforge.dev/projectforge/app/util"
)

const goStdBin, keyUser = "go", "user"

func (t *TemplateContext) Title() string {
	if t.Name != "" {
		return t.Name
	}
	return t.Key
}

func (t *TemplateContext) CleanKey() string {
	return clean(t.Key)
}

func (t *TemplateContext) NotebookPort() int {
	return t.Port + 10
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

func (t *TemplateContext) DangerousOK() bool {
	return !t.Build.SafeMode
}

func (t *TemplateContext) ModuleMarkdown() string {
	return strings.Join(util.ArraySorted(lo.Map(t.Modules, func(m string, _ int) string {
		return fmt.Sprintf("- [%s](./doc/module/%s.md)", m, m)
	})), t.Linebreak)
}

func (t *TemplateContext) PortIncremented(i int) int {
	return t.Port + i
}

func (t *TemplateContext) BuildAndroid() bool {
	ret := t.HasModules("android") && t.Build.Android
	return ret
}

func (t *TemplateContext) BuildIOS() bool {
	return t.HasModules("ios") && t.Build.IOS
}

func (t *TemplateContext) BuildDesktop() bool {
	return t.HasModules("desktop") && t.Build.Desktop
}

func (t *TemplateContext) BuildMobile() bool {
	return t.BuildIOS() || t.BuildAndroid()
}

func (t *TemplateContext) BuildWASM() bool {
	return t.HasModules("wasmserver") && t.Build.WASM
}

func (t *TemplateContext) BuildNotarize() bool {
	return t.HasModules("notarize") && t.Build.Notarize
}

func (t *TemplateContext) UsesLib() bool {
	return t.BuildMobile() || t.Build.Desktop
}

func (t *TemplateContext) HasSlack() bool {
	return t.Info.Slack != ""
}

func (t *TemplateContext) HasAccount() bool {
	return t.HasModules("oauth")
}

func (t *TemplateContext) HasUser() bool {
	return t.HasModules(keyUser)
}

func (t *TemplateContext) IsNotarized() bool {
	return t.HasModule("notarize") && t.Build != nil && t.Build.Notarize
}

func (t *TemplateContext) IsArmAndMips() bool {
	return t.Build.HasArm() && t.Build.LinuxMIPS
}

func (t *TemplateContext) DatabaseUIOpts() (bool, bool, bool) {
	cfg, _ := t.Config.GetMap("databaseui", true)
	if len(cfg) == 0 {
		return true, false, false
	}
	sqleditor := cfg.GetBoolOpt("sqleditor")
	readonly := cfg.GetBoolOpt("readonly")
	saveuser := cfg.GetBoolOpt("saveuser")
	return sqleditor, readonly, saveuser
}

func (t *TemplateContext) DatabaseUISQLEditor() bool {
	sqleditor, _, _ := t.DatabaseUIOpts()
	return sqleditor
}

func (t *TemplateContext) DatabaseUIReadOnly() bool {
	_, readonly, _ := t.DatabaseUIOpts()
	return readonly
}

func (t *TemplateContext) DatabaseUISaveUser() bool {
	_, _, saveUser := t.DatabaseUIOpts()
	return saveUser
}

func (t *TemplateContext) GoVersionSafe() string {
	return util.OrDefault(t.Info.GoVersion, DefaultGoVersion)
}

func (t *TemplateContext) GoBinarySafe() string {
	return util.OrDefault(t.Info.GoBinary, goStdBin)
}

func (t *TemplateContext) Placeholder(idx int) string {
	if t.DatabaseEngine == util.DatabasePostgreSQL || t.DatabaseEngine == util.DatabaseSQLite {
		return fmt.Sprintf("$%d", idx)
	}
	if t.DatabaseEngine == util.DatabaseSQLServer {
		return fmt.Sprintf("@p%d", idx)
	}
	return "?"
}

func (t *TemplateContext) TypeUUID() string {
	if t.DatabaseEngine == util.DatabaseSQLite {
		return "text"
	}
	return "uuid"
}

func (t *TemplateContext) HasExport() bool {
	return t.HasModules("export")
}

func (t *TemplateContext) MySQL() bool {
	return t.DatabaseEngine == util.DatabaseMySQL || t.HasModule(util.DatabaseMySQL)
}

func (t *TemplateContext) PostgreSQL() bool {
	return t.DatabaseEngine == util.DatabasePostgreSQL || t.HasModule(util.DatabasePostgreSQL)
}

func (t *TemplateContext) SQLite() bool {
	return t.DatabaseEngine == util.DatabaseSQLite || t.HasModule(util.DatabaseSQLite)
}

func (t *TemplateContext) SQLServer() bool {
	return t.DatabaseEngine == util.DatabaseSQLServer || t.HasModule(util.DatabaseSQLServer)
}

func (t *TemplateContext) SQLServerOnly() bool {
	return t.SQLServer() && !t.MySQL() && !t.PostgreSQL() && !t.SQLite()
}
