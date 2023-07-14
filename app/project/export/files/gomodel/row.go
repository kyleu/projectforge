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
	lo.ForEach(helper.ImportsForTypes("row", args.Database(), m.Columns.Types()...), func(imp *golang.Import, _ int) {
		g.AddImport(imp)
	})
	g.AddImport(helper.ImpStrings, helper.ImpAppUtil, helper.ImpFmt)
	if err := helper.SpecialImports(g, m.Columns, m.PackageWithGroup(""), args.Enums); err != nil {
		return nil, err
	}
	g.AddImport(helper.ImpLo)
	lo.ForEach(m.Columns, func(col *model.Column, _ int) {
		if col.Nullable && (col.Type.Key() == types.KeyString || col.Type.Key() == types.KeyInt || col.Type.Key() == types.KeyBool) {
			g.AddImport(helper.ImpSQL)
		}
		if col.Type.Key() == types.KeyUUID && args.Database() == model.SQLServer {
			if col.Nullable {
				g.AddImport(helper.ImpDatabase)
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
	mrow, err := modelRow(m, args.Enums, args.Database())
	if err != nil {
		return nil, err
	}
	g.AddBlocks(mrow)
	mrm, err := modelRowToModel(g, m, args.Database())
	if err != nil {
		return nil, err
	}
	g.AddBlocks(mrm, modelRowArray(), modelRowArrayTransformer(m), defaultWC(m, args.Database()))
	return g.Render(addHeader, linebreak)
}

func modelTableCols(m *model.Model) (*golang.Block, error) {
	ret := golang.NewBlock("Columns", "procedural")
	ret.W("var (")
	ret.W("\ttable         = %q", m.Name)
	ret.W("\ttableQuoted   = fmt.Sprintf(\"%%q\", table)")
	ret.W("\tcolumns       = []string{%s}", strings.Join(m.Columns.NamesQuoted(), ", "))
	ret.W("\tcolumnsQuoted = util.StringArrayQuoted(columns)")
	ret.W("\tcolumnsString = strings.Join(columnsQuoted, \", \")")
	if m.IsRevision() {
		hc := m.HistoryColumns(true)
		hcp := hc.Col.Proper()
		ret.WB()
		constCols := strings.Join(hc.Const.NamesQuoted(), ", ")
		ret.W("\tcolumns%s = util.StringArrayQuoted([]string{%s})", util.StringPad("Core", len(hcp)), constCols)
		varCols := strings.Join(hc.Var.NamesQuoted(), ", ")
		ret.W("\tcolumns%s = util.StringArrayQuoted([]string{%s})", hcp, varCols)
		ret.WB()
		ret.W("\ttable%s       = table + \"_%s\"", hcp, hc.Col.Name)
		ret.W("\ttable%sQuoted = fmt.Sprintf(\"%%%%q\", table%s)", hcp, hcp)
		joinClause := fmt.Sprintf("%%%%q %s join %%%%q %sr on ", m.FirstLetter(), m.FirstLetter())
		var joins []string
		for idx, col := range hc.Const {
			if col.PK || col.HasTag("current_revision") {
				rCol := hc.Var[idx]
				if !(rCol.PK || rCol.HasTag(model.RevisionType)) {
					return nil, errors.Errorf("invalid revision column [%s] at index [%d]", rCol.Name, idx)
				}
				joins = append(joins, fmt.Sprintf("%s.%s = %sr.%s", m.FirstLetter(), col.NameQuoted(), m.FirstLetter(), rCol.NameQuoted()))
			}
		}
		joinClause += strings.Join(joins, " and ")
		ret.W("\ttables%s = fmt.Sprintf(`%s`, table, table%s) //nolint", util.StringPad("Joined", len(hcp)+5), joinClause, hcp)
	}
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
		ret.W("\t%s %s `db:%q`", util.StringPad(c.Proper(), maxColLength), util.StringPad(gdt, maxTypeLength), c.Name)
	}
	ret.W("}")
	return ret, nil
}

func modelRowToModel(g *golang.File, m *model.Model, database string) (*golang.Block, error) {
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
			t := "any"
			if c.Type.IsListOf(types.NewString()) {
				t = "string"
			}
			ret.W("\t%sArg := []%s{}", c.Camel(), t)
			ret.W("\t_ = util.FromJSON(r.%s, &%sArg)", c.Proper(), c.Camel())
			refs = append(refs, fmt.Sprintf("%s %sArg", k, c.Camel()))
		case types.KeyMap, types.KeyValueMap:
			ret.W("\t%sArg := util.ValueMap{}", c.Camel())
			ret.W("\t_ = util.FromJSON(r.%s, &%sArg)", c.Proper(), c.Camel())
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
			case database == model.SQLServer && c.Type.Key() == types.KeyUUID:
				if c.Nullable {
					refs = append(refs, fmt.Sprintf("%s database.UUIDFromGUID(r.%s)", k, c.Proper()))
				} else {
					refs = append(refs, fmt.Sprintf("%s util.UUIDFromStringOK(r.%s.String())", k, c.Proper()))
				}
			default:
				refs = append(refs, fmt.Sprintf("%s r.%s", k, c.Proper()))
			}
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
		if database == model.SQLServer {
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
