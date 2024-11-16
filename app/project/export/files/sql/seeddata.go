package sql

import (
	"fmt"
	"slices"
	"strings"

	"github.com/pkg/errors"
	"github.com/samber/lo"

	"projectforge.dev/projectforge/app/file"
	"projectforge.dev/projectforge/app/lib/metamodel/model"
	"projectforge.dev/projectforge/app/lib/types"
	"projectforge.dev/projectforge/app/project/export/files/helper"
	"projectforge.dev/projectforge/app/project/export/golang"
	"projectforge.dev/projectforge/app/util"
)

const ind1 = "  "

func SeedData(m *model.Model, database string, linebreak string) (*file.File, error) {
	g := golang.NewGoTemplate([]string{"queries", "seeddata"}, fmt.Sprintf("seed_%s.sql", m.Name))
	seed, err := sqlSeedData(m, database)
	if err != nil {
		return nil, err
	}
	g.AddBlocks(seed)
	return g.Render(linebreak)
}

func sqlSeedData(m *model.Model, database string) (*golang.Block, error) {
	ret := golang.NewBlock("SQLCreate", "sql")
	ret.W(sqlFunc(m.Proper() + "SeedData"))
	err := sqlSeedDataColumns(m, ret, m.Table(), m.Columns, database)
	if err != nil {
		return nil, err
	}
	ret.W(sqlEnd())
	return ret, nil
}

//nolint:gocognit
func sqlSeedDataColumns(m *model.Model, block *golang.Block, tableName string, cols model.Columns, database string) error {
	block.WF("insert into %q (", tableName)
	block.W(ind1 + strings.Join(cols.NamesQuoted(), ", "))
	block.W(") values (")
	for idx, row := range m.SeedData {
		if len(row) != len(m.Columns) {
			return errors.Errorf("seed data row [%d] expected [%d] columns, but only [%d] were provided", idx+1, len(m.Columns), len(row))
		}
		var vs []string
		for _, col := range cols {
			colIdx := slices.IndexFunc(m.Columns, func(c *model.Column) bool {
				return col.Name == c.Name
			})
			if colIdx == -1 && strings.HasPrefix(col.SQL(), m.Table()+"_") {
				trimmed := strings.TrimPrefix(col.SQL(), m.Table()+"_")
				colIdx = slices.IndexFunc(m.Columns, func(c *model.Column) bool {
					return trimmed == c.Name
				})
			}
			if colIdx == -1 && strings.HasPrefix(col.SQL(), "current_") {
				trimmed := strings.TrimPrefix(col.SQL(), "current_")
				colIdx = slices.IndexFunc(m.Columns, func(c *model.Column) bool {
					return trimmed == c.Name
				})
			}
			if colIdx == -1 {
				return errors.Errorf("unable to find column [%s]", col.Name)
			}
			cell := row[colIdx]
			cellStr := fmt.Sprint(cell)
			switch col.Type.Key() {
			case types.KeyString, types.KeyEnum:
				vs = append(vs, processString(cellStr, "''"))
			case types.KeyDate, types.KeyTimestamp, types.KeyTimestampZoned:
				if cellStr == helper.TextNil {
					vs = append(vs, helper.TextNull)
				} else if _, err := util.TimeFromString(cellStr); err == nil {
					vs = append(vs, processString(cellStr, "''"))
				} else {
					vs = append(vs, cellStr)
				}
			case types.KeyUUID:
				vs = append(vs, processString(cellStr, "'00000000-0000-0000-0000-000000000000'"))
			case types.KeyList:
				vs = append(vs, processList(cell, cellStr))
			case types.KeyInt:
				if cellStr == helper.TextNil {
					vs = append(vs, "0")
					continue
				}
				vs = append(vs, fmt.Sprintf("%.0f", cell))
			case types.KeyFloat:
				if cellStr == helper.TextNil {
					vs = append(vs, "0")
					continue
				}
				vs = append(vs, fmt.Sprintf("%f", cell))
			case types.KeyMap, types.KeyValueMap, types.KeyReference, types.KeyNumeric:
				if cellStr == helper.TextNil {
					vs = append(vs, helper.TextNull)
					continue
				}
				switch cell.(type) {
				case string:
					vs = append(vs, "'"+cellStr+"'")
				default:
					vs = append(vs, "'"+strings.ReplaceAll(util.ToJSONCompact(cell), "'", "''")+"'")
				}
			default:
				if cellStr == helper.TextNil {
					vs = append(vs, helper.TextNull)
					continue
				}
				vs = append(vs, cellStr)
			}
		}
		block.W(ind1 + strings.Join(vs, ", "))
		if idx < len(m.SeedData)-1 {
			block.W("), (")
		}
	}
	if database == util.DatabaseSQLServer {
		block.W(");")
	} else {
		block.W(") on conflict do nothing;")
	}
	return nil
}

func processString(cellStr string, dflt string) string {
	if cellStr == helper.TextNil {
		return dflt
	}
	return "'" + clean(cellStr) + "'"
}

func processList(cell any, cellStr string) string {
	if cellStr == helper.TextNil {
		return "'[]'"
	}
	a, ok := cell.([]any)
	if !ok {
		s, ok := cell.([]string)
		if ok {
			a = lo.ToAnySlice(s)
		} else {
			str, ok := cell.(string)
			if ok {
				return "'" + str + "'"
			}
			return "'[\"error:invalid_type\"]'"
		}
	}
	ret := util.NewStringSlice(make([]string, 0, len(a)))
	lo.ForEach(a, func(x any, _ int) {
		switch t := x.(type) {
		case string:
			ret.Pushf("%q", t)
		default:
			ret.Push(fmt.Sprint(t))
		}
	})
	return "'[" + clean(ret.Join(", ")) + "]'"
}

var cleanRpl *strings.Replacer

func clean(s string) string {
	if cleanRpl == nil {
		cleanRpl = strings.NewReplacer("'", "''", "\f", "", "\v", "", "\u0000", "", "\u0082", "", "%", "%%", "{{{", "{ {{")
	}
	return cleanRpl.Replace(s)
}
