package sql

import (
	"github.com/kyleu/projectforge/app/export/golang"
	"github.com/kyleu/projectforge/app/export/model"
	"github.com/kyleu/projectforge/app/file"
)

func MigrationAll(models model.Models, addHeader bool) (*file.File, error) {
	g := golang.NewGoTemplate([]string{"queries", "ddl"}, "all.sql")
	g.AddBlocks(sqlDropAll(models), sqlCreateAll(models))
	return g.Render(addHeader)
}

func sqlDropAll(models model.Models) *golang.Block {
	ret := golang.NewBlock("SQLDropAll", "sql")
	ret.W("-- {%% func DropAll() %%}")
	for i := len(models) - 1; i >= 0; i-- {
		ret.W("-- {%%%%= %sDrop() %%%%}", models[i].Proper())
	}
	ret.W("-- {%% endfunc %%}")
	return ret
}

func sqlCreateAll(models model.Models) *golang.Block {
	ret := golang.NewBlock("SQLCreateAll", "sql")
	ret.W("-- {%% func CreateAll() %%}")
	for _, model := range models {
		ret.W("-- {%%%%= %sCreate() %%%%}", model.Proper())
	}
	ret.W("-- {%% endfunc %%}")
	return ret
}
