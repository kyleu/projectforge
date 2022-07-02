package sql

import (
	"fmt"
	"strings"

	"projectforge.dev/projectforge/app/file"
	"projectforge.dev/projectforge/app/project/export/golang"
	"projectforge.dev/projectforge/app/project/export/model"
	"projectforge.dev/projectforge/app/util"
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
				vs = append(vs, processString(cellStr, "''"))
			case "uuid":
				vs = append(vs, processString(cellStr, "'00000000-0000-0000-0000-000000000000'"))
			case "list":
				vs = append(vs, processList(cell, cellStr))
			case "int", "int64":
				if cellStr == nilStr {
					vs = append(vs, "0")
					continue
				}
				vs = append(vs, fmt.Sprintf("%.0f", cell))
			case "float", "float64":
				if cellStr == nilStr {
					vs = append(vs, "0")
					continue
				}
				vs = append(vs, fmt.Sprintf("%f", cell))
			case "map", "valuemap":
				if cellStr == nilStr {
					vs = append(vs, "null")
					continue
				}
				vs = append(vs, "'"+util.ToJSONCompact(cell)+"'")
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

func processString(cellStr string, dflt string) string {
	if cellStr == nilStr {
		return dflt
	}
	return "'" + clean(cellStr) + "'"
}

func processList(cell any, cellStr string) string {
	if cellStr == nilStr {
		return "'[]'"
	}
	a, ok := cell.([]any)
	if !ok {
		return "'[\"error:invalid_type\"]'"
	}
	ret := make([]string, 0, len(a))
	for _, x := range a {
		switch t := x.(type) {
		case string:
			ret = append(ret, "\""+t+"\"")
		default:
			ret = append(ret, fmt.Sprint(t))
		}
	}
	return "'[" + clean(strings.Join(ret, ", ")) + "]'"
}

var cleanRpl *strings.Replacer

func clean(s string) string {
	if cleanRpl == nil {
		cleanRpl = strings.NewReplacer("'", "''", "\f", "", "\v", "", "\u0000", "", "\u0082", "", "%", "%%", "{{{", "{ {{")
	}
	return cleanRpl.Replace(s)
}
