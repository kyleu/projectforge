package template

import (
	"fmt"
	"strings"

	"github.com/samber/lo"

	"projectforge.dev/projectforge/app/project"
	"projectforge.dev/projectforge/app/util"
)

const goStdBin, keyUser = "go", "user"

func (t *Context) Title() string {
	if t.Name != "" {
		return t.Name
	}
	return t.Key
}

func (t *Context) CleanKey() string {
	return project.CleanKey(t.Key)
}

func (t *Context) CleanName() string {
	return project.CleanKey(t.Name)
}

func (t *Context) NotebookPort() int {
	return t.Port + 10
}

func (t *Context) AdditionalPortList() []int {
	return util.ArraySorted(lo.Values(t.PortOffsets))
}

func (t *Context) KeyProper() string {
	return strings.ToUpper(t.Key[:1]) + t.Key[1:]
}

func (t *Context) NameCompressed() string {
	return strings.ReplaceAll(t.Name, " ", "")
}

func (t *Context) SourceTrimmed() string {
	return strings.TrimPrefix(strings.TrimPrefix(t.Info.Sourcecode, "http://"), "https://")
}

func (t *Context) ExecSafe() string {
	if t.Exec == t.Key {
		return "AppKey"
	}
	return fmt.Sprintf("%q", t.Exec)
}

func (t *Context) ArgsQuoted() string {
	return strings.Join(util.StringArrayQuoted(append(util.ArrayCopy(t.Args), t.Exec)), ", ")
}

func (t *Context) DangerousOK() bool {
	return !t.Build.SafeMode
}

func (t *Context) ModuleMarkdown() string {
	return util.StringJoin(util.ArraySorted(lo.Map(t.Modules, func(m string, _ int) string {
		return fmt.Sprintf("- [%s](./doc/module/%s.md)", m, m)
	})), t.Linebreak)
}

func (t *Context) PortIncremented(i int) int {
	return t.Port + i
}

func (t *Context) BuildAndroid() bool {
	ret := t.HasModules("android") && t.Build.Android
	return ret
}

func (t *Context) BuildIOS() bool {
	return t.HasModules("ios") && t.Build.IOS
}

func (t *Context) BuildDesktop() bool {
	return t.HasModules("desktop") && t.Build.Desktop
}

func (t *Context) BuildMobile() bool {
	return t.BuildIOS() || t.BuildAndroid()
}

func (t *Context) BuildWASM() bool {
	return t.HasModules("wasmserver") && t.Build.WASM
}

func (t *Context) BuildNotarize() bool {
	return t.HasModules("notarize") && t.Build.Notarize
}

func (t *Context) UsesLib() bool {
	return t.BuildMobile() || t.Build.Desktop
}

func (t *Context) HasSlack() bool {
	return t.Info.Slack != ""
}

func (t *Context) HasAccount() bool {
	return t.HasModules("oauth")
}

func (t *Context) HasUser() bool {
	return t.HasModules(keyUser)
}

func (t *Context) IsNotarized() bool {
	return t.HasModule("notarize") && t.Build != nil && t.Build.Notarize
}

func (t *Context) IsArmAndMips() bool {
	return t.Build.HasArm() && t.Build.LinuxMIPS
}

func (t *Context) DatabaseUIOpts() (bool, bool, bool) {
	cfg, _ := t.Config.GetMap("databaseui", true)
	if len(cfg) == 0 {
		return true, false, false
	}
	sqleditor := cfg.GetBoolOpt("sqleditor")
	readonly := cfg.GetBoolOpt("readonly")
	saveuser := cfg.GetBoolOpt("saveuser")
	return sqleditor, readonly, saveuser
}

func (t *Context) DatabaseUISQLEditor() bool {
	sqleditor, _, _ := t.DatabaseUIOpts()
	return sqleditor
}

func (t *Context) DatabaseUIReadOnly() bool {
	_, readonly, _ := t.DatabaseUIOpts()
	return readonly
}

func (t *Context) DatabaseUISaveUser() bool {
	_, _, saveUser := t.DatabaseUIOpts()
	return saveUser
}

func (t *Context) GoVersionSafe() string {
	return util.OrDefault(t.Info.GoVersion, project.DefaultGoVersion)
}

func (t *Context) GoMajorVersionSafe() string {
	v := util.OrDefault(t.Info.GoVersion, project.DefaultGoVersion)
	if strings.Count(v, ".") <= 1 {
		return v
	}
	return v[:strings.LastIndex(v, ".")]
}

func (t *Context) GoBinarySafe() string {
	return util.OrDefault(t.Info.GoBinary, goStdBin)
}

func (t *Context) Placeholder(idx int) string {
	if t.DatabaseEngine == util.DatabasePostgreSQL || t.DatabaseEngine == util.DatabaseSQLite {
		return fmt.Sprintf("$%d", idx)
	}
	if t.DatabaseEngine == util.DatabaseSQLServer {
		return fmt.Sprintf("@p%d", idx)
	}
	return "?"
}

func (t *Context) TypeUUID() string {
	if t.DatabaseEngine == util.DatabaseSQLite {
		return "text"
	}
	return "uuid"
}

func (t *Context) HasExport() bool {
	return t.HasModules("export") && !t.ExportArgs.Empty()
}

func (t *Context) HasExportModels() bool {
	return t.HasModules("export") && !t.ExportArgs.EmptyModels()
}

func (t *Context) HasDocker() bool {
	return !t.Build.SkipDocker
}

func (t *Context) HasDatabase() bool {
	return t.DatabaseEngine != ""
}

func (t *Context) HasDatabaseOrMetamodel() bool {
	return t.HasDatabase() || t.HasModule("metamodel")
}

func (t *Context) MySQL() bool {
	return t.DatabaseEngine == util.DatabaseMySQL || t.HasModule(util.DatabaseMySQL)
}

func (t *Context) MySQLOrExport() bool {
	return t.MySQL() || t.HasExport()
}

func (t *Context) MySQLOrMetamodel() bool {
	return t.MySQL() || t.HasModule("metamodel")
}

func (t *Context) MySQLOnly() bool {
	return t.MySQL() && !t.PostgreSQL() && !t.SQLite() && !t.SQLServer()
}

func (t *Context) PostgreSQL() bool {
	return t.DatabaseEngine == util.DatabasePostgreSQL || t.HasModule(util.DatabasePostgreSQL)
}

func (t *Context) PostgreSQLOrExport() bool {
	return t.PostgreSQL() || t.HasExport()
}

func (t *Context) PostgreSQLOrMetamodel() bool {
	return t.PostgreSQL() || t.HasModule("metamodel")
}

func (t *Context) PostgreSQLOnly() bool {
	return t.PostgreSQL() && !t.MySQL() && !t.SQLite() && !t.SQLServer()
}

func (t *Context) SQLite() bool {
	return t.DatabaseEngine == util.DatabaseSQLite || t.HasModule(util.DatabaseSQLite)
}

func (t *Context) SQLiteOrExport() bool {
	return t.SQLite() || t.HasExport()
}

func (t *Context) SQLiteOrMetamodel() bool {
	return t.SQLite() || t.HasModule("metamodel")
}

func (t *Context) SQLiteOnly() bool {
	return t.SQLite() && !t.MySQL() && !t.PostgreSQL() && !t.SQLServer()
}

func (t *Context) SQLServer() bool {
	return t.DatabaseEngine == util.DatabaseSQLServer || t.HasModule(util.DatabaseSQLServer)
}

func (t *Context) SQLServerOrExport() bool {
	return t.SQLServer() || t.HasExport()
}

func (t *Context) SQLServerOrMetamodel() bool {
	return t.SQLServer() || t.HasModule("metamodel")
}

func (t *Context) SQLServerOnly() bool {
	return t.SQLServer() && !t.MySQL() && !t.PostgreSQL() && !t.SQLite()
}
