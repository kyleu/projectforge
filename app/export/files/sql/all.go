package sql

import (
	"projectforge.dev/projectforge/app/export/golang"
	"projectforge.dev/projectforge/app/export/model"
	"projectforge.dev/projectforge/app/file"
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
	for _, m := range models {
		ret.W("-- {%%%%= %sCreate() %%%%}", m.Proper())
	}
	ret.W("-- {%% endfunc %%}")
	return ret
}

func SeedDataAll(models model.Models) (*file.File, error) {
	g := golang.NewGoTemplate([]string{"queries", "seeddata"}, "all.sql")
	g.AddBlocks(sqlSeedAll(models))
	return g.Render(false)
}

func sqlSeedAll(models model.Models) *golang.Block {
	ret := golang.NewBlock("SQLSeedDataAll", "sql")
	ret.W("-- {%% func SeedDataAll() %%}")
	for _, m := range models {
		if len(m.SeedData) > 0 {
			ret.W("-- {%%%%= %sSeedData() %%%%}", m.Proper())
		}
	}
	ret.W("-- {%% endfunc %%}")
	return ret
}
