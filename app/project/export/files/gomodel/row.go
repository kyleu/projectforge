package gomodel

import (
	"fmt"
	"strings"

	"github.com/pkg/errors"
	"github.com/samber/lo"

	"projectforge.dev/projectforge/app/file"
	"projectforge.dev/projectforge/app/lib/types"
	"projectforge.dev/projectforge/app/project/export/enum"
	"projectforge.dev/projectforge/app/project/export/files/helper"
	"projectforge.dev/projectforge/app/project/export/golang"
	"projectforge.dev/projectforge/app/project/export/model"
	"projectforge.dev/projectforge/app/util"
)

func Row(m *model.Model, args *model.Args, addHeader bool, linebreak string) (*file.File, error) {
	g := golang.NewFile(m.Package, []string{"app", m.PackageWithGroup("")}, "row")
	lo.ForEach(helper.ImportsForTypes("row", args.Database, m.Columns.Types()...), func(imp *golang.Import, _ int) {
		g.AddImport(imp)
	})
	g.AddImport(helper.ImpStrings, helper.ImpAppUtil, helper.ImpFmt)
	if err := helper.SpecialImports(g, m.Columns, m.PackageWithGroup(""), args.Enums); err != nil {
		return nil, err
	}
	g.AddImport(helper.ImpLo)
	lo.ForEach(m.Columns, func(col *model.Column, _ int) {
		if col.Nullable && (types.IsString(col.Type) || types.IsInt(col.Type) || types.IsBool(col.Type)) {
			g.AddImport(helper.ImpSQL)
		}
		if col.Type.Key() == types.KeyUUID && args.Database == util.DatabaseSQLServer {
			if col.Nullable {
				g.AddImport(helper.ImpAppDatabase)
			} else {
				g.AddImport(helper.ImpMSSQL)
			}
		}
	})
	if tc, err := modelTableCols(m); err == nil {
		g.AddBlocks(tc)
	} else {
		return nil, err
	}
	mrow, err := modelRow(m, args.Enums, args.Database)
	if err != nil {
		return nil, err
	}
	g.AddBlocks(mrow)
	mrm, err := modelRowToModel(g, m, args.Enums, args.Database)
	if err != nil {
		return nil, err
	}
	g.AddBlocks(mrm, modelRowArray(), modelRowArrayTransformer(m), defaultWC(m, args.Database))
	return g.Render(addHeader, linebreak)
}

func modelTableCols(m *model.Model) (*golang.Block, error) {
	ret := golang.NewBlock("Columns", "procedural")
	ret.SkipDecl = true
	ret.W("var (")
	ret.W("\ttable         = %q", m.Name)
	ret.W("\ttableQuoted   = fmt.Sprintf(\"%%q\", table)")
	cols := fmt.Sprintf("\tcolumns       = []string{%s}", strings.Join(m.Columns.NamesQuoted(), ", "))
	if len(cols) > 160 {
		cols += " //nolint:lll"
	}
	ret.W(cols)
	ret.W("\tcolumnsQuoted = util.StringArrayQuoted(columns)")
	ret.W("\tcolumnsString = strings.Join(columnsQuoted, \", \")")
	ret.W(")")
	return ret, nil
}

func modelRow(m *model.Model, enums enum.Enums, database string) (*golang.Block, error) {
	ret := golang.NewBlock(m.Proper()+"Row", "struct")
	ret.W("type row struct {")
	maxColLength := m.Columns.MaxCamelLength()
	maxTypeLength := m.Columns.MaxGoRowTypeLength(m.Package, enums, database)
	for _, c := range m.Columns {
		gdt, err := c.ToGoRowType(m.Package, enums, database)
		if err != nil {
			return nil, err
		}
		ret.W("\t%s %s `db:%q json:%q`", util.StringPad(c.Proper(), maxColLength), util.StringPad(gdt, maxTypeLength), c.Name, c.Name)
	}
	ret.W("}")
	return ret, nil
}

func modelRowToModel(g *golang.File, m *model.Model, enums enum.Enums, database string) (*golang.Block, error) {
	ret := golang.NewBlock(m.Proper(), "func")
	ret.W("func (r *row) To%s() *%s {", m.Proper(), m.Proper())
	ret.W("\tif r == nil {")
	ret.W("\t\treturn nil")
	ret.W("\t}")
	refs := make([]string, 0, len(m.Columns))
	pad := m.Columns.MaxCamelLength() + 1
	for _, c := range m.Columns {
		k := util.StringPad(c.Proper()+":", pad)
		switch c.Type.Key() {
		case types.KeyAny:
			ret.W("\tvar %sArg any", c.Camel())
			ret.W("\t_ = util.FromJSON(r.%s, &%sArg)", c.Proper(), c.Camel())
			refs = append(refs, fmt.Sprintf("%s %sArg", k, c.Camel()))
		case types.KeyList:
			t := "[]any"
			decoder := "r." + c.Proper()
			switch c.Type.ListType().Key() {
			case types.KeyString:
				t = "[]string"
				if database == util.DatabaseSQLite {
					decoder = "[]byte(" + decoder + ")"
				}
			case types.KeyInt:
				t = "[]int"
				if database == util.DatabaseSQLite {
					decoder = "[]byte(" + decoder + ")"
				}
			case types.KeyEnum:
				if e, _ := model.AsEnumInstance(c.Type.ListType(), enums); e != nil {
					t = e.ProperPlural()
					decoder = "[]byte(" + decoder + ")"
				}
			}
			ret.W("\t%sArg := %s{}", c.Camel(), t)
			ret.W("\t_ = util.FromJSON(%s, &%sArg)", decoder, c.Camel())
			refs = append(refs, fmt.Sprintf("%s %sArg", k, c.Camel()))
		case types.KeyMap, types.KeyValueMap:
			decoder := "r." + c.Proper()
			if database == util.DatabaseSQLite {
				decoder = "[]byte(" + decoder + ")"
			}
			ret.W("\t%sArg, _ := util.FromJSONMap(%s)", c.Camel(), decoder)
			refs = append(refs, fmt.Sprintf("%s %sArg", k, c.Camel()))
		case types.KeyReference:
			ref, err := model.AsRef(c.Type)
			if err != nil {
				return nil, errors.Wrap(err, "invalid ref")
			}
			if ref.Pkg.Last() == m.Package {
				ret.W("\t%sArg := &%s{}", c.Camel(), ref.K)
			} else {
				ret.W("\t%sArg := &%s.%s{}", c.Camel(), ref.Pkg.Last(), ref.K)
				g.AddImport(golang.NewImport(golang.ImportTypeApp, ref.Pkg.ToPath()))
			}
			ret.W("\t_ = util.FromJSON(r.%s, %sArg)", c.Proper(), c.Camel())
			refs = append(refs, fmt.Sprintf("%s %sArg", k, c.Camel()))
		default:
			refs = append(refs, defaultRowToModel(k, c, database)...)
		}
	}
	ret.W("\treturn &%s{", m.Proper())
	lo.ForEach(refs, func(ref string, _ int) {
		ret.W("\t\t%s,", ref)
	})
	ret.W("\t}")
	ret.W("}")
	return ret, nil
}

func defaultRowToModel(k string, c *model.Column, database string) []string {
	var refs []string
	switch {
	case c.Type.Scalar() && c.Nullable:
		switch c.Type.Key() {
		case types.KeyString:
			refs = append(refs, fmt.Sprintf("%s r.%s.String", k, c.Proper()))
		case types.KeyInt:
			refs = append(refs, fmt.Sprintf("%s int(r.%s.Int64)", k, c.Proper()))
		case types.KeyFloat:
			refs = append(refs, fmt.Sprintf("%s r.%s.Float64", k, c.Proper()))
		case types.KeyBool:
			refs = append(refs, fmt.Sprintf("%s r.%s.Bool", k, c.Proper()))
		default:
			refs = append(refs, fmt.Sprintf("%s r.%s", k, c.Proper()))
		}
	case database == util.DatabaseSQLServer && c.Type.Key() == types.KeyUUID:
		if c.Nullable {
			refs = append(refs, fmt.Sprintf("%s database.UUIDFromGUID(r.%s)", k, c.Proper()))
		} else {
			refs = append(refs, fmt.Sprintf("%s util.UUIDFromStringOK(r.%s.String())", k, c.Proper()))
		}
	default:
		refs = append(refs, fmt.Sprintf("%s r.%s", k, c.Proper()))
	}
	return refs
}

func modelRowArray() *golang.Block {
	ret := golang.NewBlock("RowArray", "type")
	ret.W("type rows []*row")
	return ret
}

func modelRowArrayTransformer(m *model.Model) *golang.Block {
	ret := golang.NewBlock(fmt.Sprintf("RowTo%s", m.ProperPlural()), "type")
	ret.W("func (x rows) To%s() %s {", m.ProperPlural(), m.ProperPlural())
	ret.W("\treturn lo.Map(x, func(d *row, _ int) *%s {", m.Proper())
	ret.W("\t\treturn d.To%s()", m.Proper())
	ret.W("\t})")
	ret.W("}")
	return ret
}

func defaultWC(m *model.Model, database string) *golang.Block {
	ret := golang.NewBlock("Columns", "procedural")
	ret.W("func defaultWC(idx int) string {")
	c := m.PKs()
	wc := make([]string, 0, len(c))
	idxs := make([]string, 0, len(c))
	lo.ForEach(c, func(col *model.Column, idx int) {
		if database == util.DatabaseSQLServer {
			wc = append(wc, fmt.Sprintf("%q = @p%%%%d", col.Name))
		} else {
			wc = append(wc, fmt.Sprintf("%q = $%%%%d", col.Name))
		}
		idxs = append(idxs, fmt.Sprintf("idx+%d", idx+1))
	})
	ret.W("\treturn fmt.Sprintf(%q, %s)", strings.Join(wc, " and "), strings.Join(idxs, ", "))
	ret.W("}")
	return ret
}
