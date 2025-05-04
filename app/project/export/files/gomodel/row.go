package gomodel

import (
	"fmt"

	"github.com/pkg/errors"
	"github.com/samber/lo"

	"projectforge.dev/projectforge/app/file"
	"projectforge.dev/projectforge/app/lib/metamodel"
	"projectforge.dev/projectforge/app/lib/metamodel/enum"
	"projectforge.dev/projectforge/app/lib/metamodel/model"
	"projectforge.dev/projectforge/app/lib/types"
	"projectforge.dev/projectforge/app/project/export/files/helper"
	"projectforge.dev/projectforge/app/project/export/golang"
	"projectforge.dev/projectforge/app/util"
)

const rDot = "r."

func Row(m *model.Model, args *metamodel.Args, linebreak string) (*file.File, error) {
	g := golang.NewFile(m.Package, []string{"app", m.PackageWithGroup("")}, "row")
	lo.ForEach(helper.ImportsForTypes("row", args.Database, m.Columns.Types()...), func(imp *model.Import, _ int) {
		g.AddImport(imp)
	})
	g.AddImport(helper.ImpAppUtil, helper.ImpFmt)
	imps, err := helper.SpecialImports(m.Columns, m.PackageWithGroup(""), args.Models, args.Enums, args.ExtraTypes)
	if err != nil {
		return nil, err
	}
	g.AddImport(imps...)
	g.AddImport(helper.ImpLo)
	imps, err = helper.EnumImports(m.Columns.Types(), m.PackageWithGroup(""), args.Enums)
	if err != nil {
		return nil, err
	}
	g.AddImport(imps...)
	g.AddImport(m.Imports.Supporting("row")...)
	lo.ForEach(m.Columns, func(col *model.Column, _ int) {
		if col.Type.Key() == types.KeyUUID && args.Database == util.DatabaseSQLServer {
			if col.Nullable {
				g.AddImport(helper.ImpAppDatabase)
			} else {
				g.AddImport(helper.ImpMSSQL)
			}
		}
	})
	if tc, err2 := modelTableCols(m); err2 == nil {
		g.AddBlocks(tc)
	} else {
		return nil, err2
	}
	mrow, err := modelRow(m, args.Enums, args.Database)
	if err != nil {
		return nil, err
	}
	g.AddBlocks(mrow)
	mrm, err := modelRowToModel(m, args.Models, args.Enums, args.ExtraTypes, args.Database)
	if err != nil {
		return nil, err
	}
	g.AddBlocks(mrm, modelRowArray(), modelRowArrayTransformer(m), defaultWC(m, args.Database))
	return g.Render(linebreak)
}

func modelTableCols(m *model.Model) (*golang.Block, error) {
	ret := golang.NewBlock("Columns", "procedural")
	ret.SkipDecl = true
	ret.W("var (")
	ret.WF("\ttable         = %q", m.Table())
	if m.Schema == "" {
		ret.W("\ttableQuoted   = fmt.Sprintf(\"%%q\", table)")
	} else {
		ret.WF("\ttableQuoted   = fmt.Sprintf(\"%%%%q.%%%%q\", %q, table)", m.Schema)
	}
	cols := fmt.Sprintf("\tcolumns       = []string{%s}", util.StringJoin(m.Columns.NotDerived().SQLQuoted(), ", "))
	if len(cols) > 160 {
		cols += " //nolint:lll"
	}
	ret.W(cols)
	ret.W("\tcolumnsQuoted = util.StringArrayQuoted(columns)")
	ret.W("\tcolumnsString = util.StringJoin(columnsQuoted, \", \")")
	ret.W(")")
	return ret, nil
}

func modelRow(m *model.Model, enums enum.Enums, database string) (*golang.Block, error) {
	ret := golang.NewBlock(m.Proper()+"Row", "struct")
	ret.W("type row struct {")
	cols := m.Columns.NotDerived()
	maxColLength := cols.MaxCamelLength()
	maxTypeLength := cols.MaxGoRowTypeLength(m.Package, enums, database)
	for _, c := range cols {
		gdt, err := c.ToGoRowType(m.Package, enums, database)
		if err != nil {
			return nil, err
		}
		if gdt == "*"+types.KeyAny {
			gdt = types.KeyAny
		}
		ret.WF("\t%s %s `db:%q json:%q`", util.StringPad(c.Proper(), maxColLength), util.StringPad(gdt, maxTypeLength), c.SQL(), c.Name)
	}
	ret.W("}")
	return ret, nil
}

//nolint:gocognit
func modelRowToModel(m *model.Model, models model.Models, enums enum.Enums, extraTypes model.Models, database string) (*golang.Block, error) {
	ba := func(decoder string) string {
		return "[]byte(" + decoder + ")"
	}
	ret := golang.NewBlock(m.Proper(), "func")
	ret.WF("func (r *row) To%s() *%s {", m.Proper(), m.Proper())
	ret.W("\tif r == nil {")
	ret.W("\t\treturn nil")
	ret.W("\t}")
	cols := m.Columns.NotDerived()
	refs := make([]string, 0, len(cols))
	pad := cols.MaxCamelLength() + 1
	for _, c := range cols {
		k := util.StringPad(c.Proper()+":", pad)
		switch c.Type.Key() {
		case types.KeyAny:
			ret.WF("\tvar %sArg any", c.Camel())
			if helper.SimpleJSON(database) {
				ret.WF("\t_ = util.FromJSON([]byte(r.%s), &%sArg)", c.Proper(), c.Camel())
			} else {
				ret.WF("\t_ = util.FromJSON(r.%s, &%sArg)", c.Proper(), c.Camel())
			}
			refs = append(refs, fmt.Sprintf("%s %sArg", k, c.Camel()))
		case types.KeyJSON:
			ret.WF("\tvar %sArg any", c.Camel())
			if helper.SimpleJSON(database) {
				ret.WF("\t_ = util.FromJSON([]byte(r.%s.V), &%sArg)", c.Proper(), c.Camel())
			} else {
				ret.WF("\t_ = util.FromJSON(r.%s.V, &%sArg)", c.Proper(), c.Camel())
			}
			refs = append(refs, fmt.Sprintf("%s %sArg", k, c.Camel()))
		case types.KeyList:
			t := "[]any"
			decoder := rDot + c.Proper()
			switch c.Type.ListType().Key() {
			case types.KeyString:
				t = "[]string"
				if helper.SimpleJSON(database) {
					decoder = ba(decoder)
				}
			case types.KeyInt:
				t = "[]int"
				if helper.SimpleJSON(database) {
					decoder = ba(decoder)
				}
			case types.KeyMap, types.KeyValueMap:
				t = "[]util.ValueMap"
				if helper.SimpleJSON(database) {
					decoder = ba(decoder)
				}
			case types.KeyOrderedMap:
				t = "[]util.OrderedMap[any]"
				if helper.SimpleJSON(database) {
					decoder = ba(decoder)
				}
			case types.KeyEnum:
				if e, _ := model.AsEnumInstance(c.Type.ListType(), enums); e != nil {
					t = e.ProperPlural()
					if e.PackageWithGroup("") != m.PackageWithGroup("") {
						t = e.Package + "." + t
					}
					decoder = ba(decoder)
				}
			}
			ret.WF("\tvar %sArg %s", c.Camel(), t)
			ret.WF("\t_ = util.FromJSON(%s, &%sArg)", decoder, c.Camel())
			refs = append(refs, fmt.Sprintf("%s %sArg", k, c.Camel()))
		case types.KeyMap, types.KeyValueMap:
			decoder := rDot + c.Proper()
			if helper.SimpleJSON(database) {
				decoder = ba(decoder)
			}
			ret.WF("\t%sArg, _ := util.FromJSONMap(%s)", c.Camel(), decoder)
			refs = append(refs, fmt.Sprintf("%s %sArg", k, c.Camel()))
		case types.KeyReference:
			decoder := rDot + c.Proper()
			if helper.SimpleJSON(database) {
				decoder = ba(decoder)
			}
			ref, _, err := helper.LoadRef(c, models, extraTypes)
			if err != nil {
				return nil, errors.Wrap(err, "invalid ref")
			}
			ret.WF("\t%sArg := %s{}", c.Camel(), ref.LastAddr(ref.Pkg.Last() != m.Package))
			ret.WF("\t_ = util.FromJSON(%s, %sArg)", decoder, c.Camel())
			refs = append(refs, fmt.Sprintf("%s %sArg", k, c.Camel()))
		default:
			refs = append(refs, defaultRowToModel(k, c, database)...)
		}
	}
	ret.WF("\treturn &%s{", m.Proper())
	lo.ForEach(refs, func(ref string, _ int) {
		ret.WF("\t\t%s,", ref)
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
			if i := types.TypeAs[*types.Int](c.Type); i != nil {
				b := util.Choose(i.Bits == 0, 64, i.Bits)
				refs = append(refs, fmt.Sprintf("%s r.%s.Int%d", k, c.Proper(), b))
			} else {
				refs = append(refs, fmt.Sprintf("%s int(r.%s.Int64)", k, c.Proper()))
			}
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
	ret.WF("func (x rows) To%s() %s {", m.ProperPlural(), m.ProperPlural())
	ret.WF("\treturn lo.Map(x, func(d *row, _ int) *%s {", m.Proper())
	ret.WF("\t\treturn d.To%s()", m.Proper())
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
			wc = append(wc, fmt.Sprintf("%q = @p%%%%d", col.SQL()))
		} else {
			wc = append(wc, fmt.Sprintf("%q = $%%%%d", col.SQL()))
		}
		idxs = append(idxs, fmt.Sprintf("idx+%d", idx+1))
	})
	ret.WF("\treturn fmt.Sprintf(%q, %s)", util.StringJoin(wc, " and "), util.StringJoin(idxs, ", "))
	ret.W("}")
	return ret
}
