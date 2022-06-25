package sql

import (
	"fmt"
	"strings"

	"projectforge.dev/projectforge/app/export/golang"
	"projectforge.dev/projectforge/app/export/model"
	"projectforge.dev/projectforge/app/file"
)

const nilStr = "<nil>"

func SeedData(m *model.Model, args *model.Args) (*file.File, error) {
	g := golang.NewGoTemplate([]string{"queries", "seeddata"}, m.Package+".sql")
	g.AddBlocks(sqlSeedData(m, args.Modules))
	return g.Render(false)
}

func sqlSeedData(m *model.Model, modules []string) *golang.Block {
	ret := golang.NewBlock("SQLCreate", "sql")
	ret.W("-- {%% func " + m.Proper() + "SeedData() %%}")
	ret.W("insert into %q (", m.Name)
	ret.W("  " + strings.Join(m.Columns.NamesQuoted(), ", "))
	ret.W(") values (")
	for idx, row := range m.SeedData {
		var vs []string
		for colIdx, col := range m.Columns {
			cell := row[colIdx]
			cellStr := fmt.Sprint(cell)
			switch col.Type.Key() {
			case "string":
				if cellStr == nilStr {
					vs = append(vs, "''")
					continue
				}
				vs = append(vs, "'"+clean(cellStr)+"'")
			case "uuid":
				if cellStr == nilStr {
					vs = append(vs, "'00000000-0000-0000-0000-000000000000'")
					continue
				}
				vs = append(vs, "'"+cellStr+"'")
			case "list":
				if cellStr == nilStr {
					vs = append(vs, "'[]'")
					continue
				}
				vs = append(vs, "'"+clean(cellStr)+"'")
			case "int", "int64", "float", "float64":
				if cellStr == nilStr {
					vs = append(vs, "0")
					continue
				}
				vs = append(vs, cellStr)
			default:
				if cellStr == nilStr {
					vs = append(vs, "null")
					continue
				}
				vs = append(vs, cellStr)
			}
		}
		ret.W("  " + strings.Join(vs, ", "))
		if idx < len(m.SeedData)-1 {
			ret.W("), (")
		}
	}
	ret.W(");")
	ret.W("-- {%% endfunc %%}")
	return ret
}

func clean(s string) string {
	return strings.ReplaceAll(strings.ReplaceAll(s, "%", "%%"), "{{{", "{ {{")
}
