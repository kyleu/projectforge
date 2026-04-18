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
