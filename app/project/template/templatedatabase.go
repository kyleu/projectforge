package template

import (
	"fmt"
	"strings"

	"github.com/samber/lo"

	"projectforge.dev/projectforge/app/util"
)

/*
	DatabaseMySQL      = "mysql"
	DatabasePostgreSQL = "postgres"
	DatabaseSQLite     = "sqlite"
	DatabaseSQLServer  = "sqlserver"
*/

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

func (t *Context) DatabaseKeys() string {
	keys := []string{}
	if t.HasModule("metamodel") {
		keys = []string{"MySQL", "PostgreSQL", "SQLite", "SQLServer"}
	} else {
		add := func(db string) {
			keys = append(keys, db)
		}
		if t.MySQL() {
			add("MySQL")
		}
		if t.PostgreSQL() {
			add("PostgreSQL")
		}
		if t.SQLite() {
			add("SQLite")
		}
		if t.SQLServer() {
			add("SQLServer")
		}
		if len(keys) == 0 {
			return ""
		}
	}
	l := util.StringArrayMaxLength(keys)
	lines := lo.Map(keys, func(x string, _ int) string {
		return fmt.Sprintf("\tDatabase%s = %q", util.StringPad(x, l), util.Choose(x == "PostgreSQL", "postgres", strings.ToLower(x)))
	})
	return "\n\n" + util.StringJoin(lines, "\n") + "\n"
}
