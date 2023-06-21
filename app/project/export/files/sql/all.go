package sql

import (
	"github.com/samber/lo"
	"projectforge.dev/projectforge/app/file"
	"projectforge.dev/projectforge/app/project/export/enum"
	"projectforge.dev/projectforge/app/project/export/golang"
	"projectforge.dev/projectforge/app/project/export/model"
)

func MigrationAll(models model.Models, enums enum.Enums, addHeader bool) (*file.File, error) {
	g := golang.NewGoTemplate([]string{"queries", "ddl"}, "all.sql")
	g.AddBlocks(sqlDropAll(models, enums), sqlCreateAll(models, enums))
	return g.Render(addHeader)
}

func sqlDropAll(models model.Models, enums enum.Enums) *golang.Block {
	ret := golang.NewBlock("SQLDropAll", "sql")
	ret.W("-- {%% func DropAll() %%}")
	for i := len(models) - 1; i >= 0; i-- {
		ret.W("-- {%%%%= %sDrop() %%%%}", models[i].Proper())
	}
	if len(enums) > 1 {
		ret.W("-- {%%= TypesDrop() %%}")
	}
	ret.W("-- {%% endfunc %%}")
	return ret
}

func sqlCreateAll(models model.Models, enums enum.Enums) *golang.Block {
	ret := golang.NewBlock("SQLCreateAll", "sql")
	ret.W("-- {%% func CreateAll() %%}")
	if len(enums) > 0 {
		ret.W("-- {%%= TypesCreate() %%}")
	}
	lo.ForEach(models, func(m *model.Model, _ int) {
		ret.W("-- {%%%%= %sCreate() %%%%}", m.Proper())
	})
	ret.W("-- {%% endfunc %%}")
	return ret
}

func SeedDataAll(models model.Models) (*file.File, error) {
	ordered := models.Sorted()
	g := golang.NewGoTemplate([]string{"queries", "seeddata"}, "all.sql")
	g.AddBlocks(sqlSeedAll(ordered))
	return g.Render(false)
}

func sqlSeedAll(models model.Models) *golang.Block {
	ret := golang.NewBlock("SQLSeedDataAll", "sql")
	ret.W("-- {%% func SeedDataAll() %%}")
	lo.ForEach(models, func(m *model.Model, index int) {
		if len(m.SeedData) > 0 {
			ret.W("-- {%%%%= %sSeedData() %%%%}", m.Proper())
		}
	})
	ret.W("-- {%% endfunc %%}")
	return ret
}
