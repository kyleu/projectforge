package gomodel

import (
	"fmt"
	"strings"

	"github.com/pkg/errors"

	"projectforge.dev/projectforge/app/file"
	"projectforge.dev/projectforge/app/lib/types"
	"projectforge.dev/projectforge/app/project/export/files/helper"
	"projectforge.dev/projectforge/app/project/export/golang"
	"projectforge.dev/projectforge/app/project/export/model"
	"projectforge.dev/projectforge/app/util"
)

func DTO(m *model.Model, args *model.Args, addHeader bool) (*file.File, error) {
	g := golang.NewFile(m.Package, []string{"app", m.PackageWithGroup("")}, "dto")
	for _, imp := range helper.ImportsForTypes("dto", m.Columns.Types()...) {
		g.AddImport(imp)
	}
	g.AddImport(helper.ImpStrings, helper.ImpAppUtil, helper.ImpFmt)
	if tc, err := modelTableCols(m, g); err == nil {
		g.AddBlocks(tc)
	} else {
		return nil, err
	}
	g.AddBlocks(modelDTO(m))
	mdm, err := modelDTOToModel(g, m)
	if err != nil {
		return nil, err
	}
	g.AddBlocks(mdm, modelDTOArray(), modelDTOArrayTransformer(m), defaultWC(m))
	return g.Render(addHeader)
}

func modelTableCols(m *model.Model, g *golang.File) (*golang.Block, error) {
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
		ret.W("")
		constCols := strings.Join(hc.Const.NamesQuoted(), ", ")
		ret.W("\tcolumns%s = util.StringArrayQuoted([]string{%s})", util.StringPad("Core", len(hcp)), constCols)
		varCols := strings.Join(hc.Var.NamesQuoted(), ", ")
		ret.W("\tcolumns%s = util.StringArrayQuoted([]string{%s})", hcp, varCols)
		ret.W("")
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
		ret.W("\ttables%s = fmt.Sprintf(`%s`, table, table%s) // nolint", util.StringPad("Joined", len(hcp)+5), joinClause, hcp)
	}
	ret.W(")")
	return ret, nil
}

func modelDTO(m *model.Model) *golang.Block {
	ret := golang.NewBlock(m.Proper()+"DTO", "struct")
	ret.W("type dto struct {")
	maxColLength := util.StringArrayMaxLength(m.Columns.CamelNames())
	maxTypeLength := m.Columns.MaxGoDTOKeyLength(m.Package)
	for _, c := range m.Columns {
		ret.W("\t%s %s `db:%q`", util.StringPad(c.Proper(), maxColLength), util.StringPad(c.ToGoDTOType(m.Package), maxTypeLength), c.Name)
	}
	ret.W("}")
	return ret
}

func modelDTOToModel(g *golang.File, m *model.Model) (*golang.Block, error) {
	ret := golang.NewBlock(m.Proper(), "func")
	ret.W("func (d *dto) To%s() *%s {", m.Proper(), m.Proper())
	ret.W("\tif d == nil {")
	ret.W("\t\treturn nil")
	ret.W("\t}")
	refs := make([]string, 0, len(m.Columns))
	pad := m.Columns.MaxCamelLength() + 1
	for _, c := range m.Columns {
		k := util.StringPad(c.Proper()+":", pad)
		switch c.Type.Key() {
		case types.KeyAny:
			ret.W("\tvar %sArg any", c.Camel())
			ret.W("\t_ = util.FromJSON(d.%s, &%sArg)", c.Proper(), c.Camel())
			refs = append(refs, fmt.Sprintf("%s %sArg", k, c.Camel()))
		case types.KeyList:
			ret.W("\t%sArg := []any{}", c.Camel())
			ret.W("\t_ = util.FromJSON(d.%s, &%sArg)", c.Proper(), c.Camel())
			refs = append(refs, fmt.Sprintf("%s %sArg", k, c.Camel()))
		case types.KeyMap, types.KeyValueMap:
			ret.W("\t%sArg := util.ValueMap{}", c.Camel())
			ret.W("\t_ = util.FromJSON(d.%s, &%sArg)", c.Proper(), c.Camel())
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
			ret.W("\t_ = util.FromJSON(d.%s, %sArg)", c.Proper(), c.Camel())
			refs = append(refs, fmt.Sprintf("%s %sArg", k, c.Camel()))
		default:
			refs = append(refs, fmt.Sprintf("%s d.%s", k, c.Proper()))
		}
	}
	ret.W("\treturn &%s{", m.Proper())
	for _, ref := range refs {
		ret.W("\t\t%s,", ref)
	}
	ret.W("\t}")
	ret.W("}")
	return ret, nil
}

func modelDTOArray() *golang.Block {
	ret := golang.NewBlock("DTOArray", "type")
	ret.W("type dtos []*dto")
	return ret
}

func modelDTOArrayTransformer(m *model.Model) *golang.Block {
	ret := golang.NewBlock(fmt.Sprintf("DTOTo%s", m.ProperPlural()), "type")
	ret.W("func (x dtos) To%s() %s {", m.ProperPlural(), m.ProperPlural())
	ret.W("\tret := make(%s, 0, len(x))", m.ProperPlural())
	ret.W("\tfor _, d := range x {")
	ret.W("\t\tret = append(ret, d.To%s())", m.Proper())
	ret.W("\t}")
	ret.W("\treturn ret")
	ret.W("}")
	return ret
}

func defaultWC(m *model.Model) *golang.Block {
	ret := golang.NewBlock("Columns", "procedural")
	ret.W("func defaultWC(idx int) string {")
	c := m.PKs()
	wc := make([]string, 0, len(c))
	idxs := make([]string, 0, len(c))
	for idx, col := range c {
		wc = append(wc, fmt.Sprintf("%q = $%%%%d", col.Name))
		idxs = append(idxs, fmt.Sprintf("idx+%d", idx+1))
	}
	ret.W("\treturn fmt.Sprintf(%q, %s)", strings.Join(wc, " and "), strings.Join(idxs, ", "))
	ret.W("}")
	return ret
}
