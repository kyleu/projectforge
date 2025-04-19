package sql

import (
	"github.com/samber/lo"

	"projectforge.dev/projectforge/app/file"
	"projectforge.dev/projectforge/app/lib/metamodel/enum"
	"projectforge.dev/projectforge/app/lib/metamodel/model"
	"projectforge.dev/projectforge/app/project/export/files/helper"
	"projectforge.dev/projectforge/app/project/export/golang"
)

func MigrationAll(models model.Models, enums enum.Enums, audit bool, linebreak string) (*file.File, error) {
	g := golang.NewGoTemplate([]string{"queries", "ddl"}, "all.sql")
	g.AddBlocks(sqlDropAll(models, enums, audit), sqlCreateAll(models, enums, audit))
	return g.Render(linebreak)
}

func sqlDropAll(models model.Models, enums enum.Enums, audit bool) *golang.Block {
	ret := golang.NewBlock("SQLDropAll", "sql")
	ret.W(sqlFunc("DropAll"))
	for i := len(models) - 1; i >= 0; i-- {
		ret.W(sqlCall(models[i].Proper() + helper.TextDrop))
	}
	if audit {
		ret.W(sqlCall("AuditDrop"))
	}
	if len(enums) > 1 {
		ret.W(sqlCall("TypesDrop"))
	}
	ret.W(sqlEnd())
	return ret
}

func sqlCreateAll(models model.Models, enums enum.Enums, audit bool) *golang.Block {
	ret := golang.NewBlock("SQLCreateAll", "sql")
	ret.W(sqlFunc("CreateAll"))
	if len(enums) > 0 {
		ret.W(sqlCall("TypesCreate"))
	}
	if audit {
		ret.W(sqlCall("AuditCreate"))
	}
	lo.ForEach(models, func(m *model.Model, _ int) {
		ret.WF(helper.TextSQLComment+"{%%%%= %sCreate() %%%%}", m.Proper())
	})
	ret.W(sqlEnd())
	return ret
}

func SeedDataAll(models model.Models, linebreak string) (*file.File, error) {
	ordered := models.Sorted()
	g := golang.NewGoTemplate([]string{"queries", "seeddata"}, "all.sql")
	g.AddBlocks(sqlSeedAll(ordered))
	return g.Render(linebreak)
}

func sqlSeedAll(models model.Models) *golang.Block {
	ret := golang.NewBlock("SQLSeedDataAll", "sql")
	ret.W(sqlFunc("SeedDataAll"))
	lo.ForEach(models.WithoutTag("external"), func(m *model.Model, _ int) {
		if len(m.SeedData) > 0 {
			ret.WF(helper.TextSQLComment+"{%%%%= %sSeedData() %%%%}", m.Proper())
		}
	})
	ret.W(sqlEnd())
	return ret
}
