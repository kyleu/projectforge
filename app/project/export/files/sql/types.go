package sql

import (
	"fmt"
	"strings"

	"github.com/samber/lo"

	"projectforge.dev/projectforge/app/file"
	"projectforge.dev/projectforge/app/lib/metamodel/enum"
	"projectforge.dev/projectforge/app/project/export/files/helper"
	"projectforge.dev/projectforge/app/project/export/golang"
	"projectforge.dev/projectforge/app/util"
)

func Types(enums enum.Enums, linebreak string, database string) (*file.File, error) {
	g := golang.NewGoTemplate([]string{"queries", "ddl"}, "types.sql")
	g.AddBlocks(typesDrop(enums, database), typesCreate(enums, database))
	return g.Render(linebreak)
}

func typesDrop(enums enum.Enums, database string) *golang.Block {
	ret := golang.NewBlock("TypesDrop", "sql")
	ret.W(sqlFunc("TypesDrop"))
	for i := len(enums) - 1; i >= 0; i-- {
		if database != util.DatabaseSQLite && database != util.DatabaseSQLServer {
			ret.WF("drop type if exists %q;", enums[i].Name)
		}
	}
	ret.W(sqlEnd())
	return ret
}

func typesCreate(enums enum.Enums, database string) *golang.Block {
	ret := golang.NewBlock("SQLCreateAll", "sql")
	ret.W(sqlFunc("TypesCreate"))
	lo.ForEach(enums, func(e *enum.Enum, _ int) {
		// create type model_service as enum ('team', 'sprint', 'estimate', 'standup', 'retro', 'story', 'feedback', 'report');
		q := make([]string, 0, len(e.Values))
		lo.ForEach(e.Values, func(x *enum.Value, _ int) {
			q = append(q, fmt.Sprintf("'%s'", strings.ReplaceAll(x.Key, "'", "''")))
		})
		switch {
		case database == util.DatabaseSQLite:
			ret.WF(helper.TextSQLComment+"skipping definition of enum [%s], since SQLite does not support custom types", e.Name)
		case database == util.DatabaseSQLServer:
			ret.WF(helper.TextSQLComment+"skipping definition of enum [%s], since SQL Server does not support custom types", e.Name)
		default:
			ret.W("do $$ begin")
			ret.WF("  create type %q as enum (%s);", e.Name, strings.Join(q, ", "))
			ret.W("exception")
			ret.W("  when duplicate_object then null;")
			ret.W("end $$;")
		}
	})
	ret.W(sqlEnd())
	return ret
}
