package sql

import (
	"fmt"
	"strings"

	"projectforge.dev/projectforge/app/file"
	"projectforge.dev/projectforge/app/project/export/enum"
	"projectforge.dev/projectforge/app/project/export/golang"
)

func Types(enums enum.Enums, addHeader bool) (*file.File, error) {
	g := golang.NewGoTemplate([]string{"queries", "ddl"}, "types.sql")
	g.AddBlocks(typesDrop(enums), typesCreate(enums))
	return g.Render(addHeader)
}

func typesDrop(enums enum.Enums) *golang.Block {
	ret := golang.NewBlock("TypesDrop", "sql")
	ret.W("-- {%% func TypesDrop() %%}")
	for i := len(enums) - 1; i >= 0; i-- {
		ret.W("drop type if exists %q;", enums[i].Name)
	}
	ret.W("-- {%% endfunc %%}")
	return ret
}

func typesCreate(enums enum.Enums) *golang.Block {
	ret := golang.NewBlock("SQLCreateAll", "sql")
	ret.W("-- {%% func TypesCreate() %%}")
	for _, e := range enums {
		// create type model_service as enum ('team', 'sprint', 'estimate', 'standup', 'retro', 'story', 'feedback', 'report');
		q := make([]string, 0, len(e.Values))
		for _, x := range e.Values {
			q = append(q, fmt.Sprintf("'%s'", strings.ReplaceAll(x, "'", "''")))
		}
		ret.W("create type %q as enum (%s);", e.Name, strings.Join(q, ", "))
	}
	ret.W("-- {%% endfunc %%}")
	return ret
}
