package gomodel

import (
	"fmt"
	"strings"

	"projectforge.dev/projectforge/app/export/files/helper"
	"projectforge.dev/projectforge/app/export/golang"
	"projectforge.dev/projectforge/app/export/model"
	"projectforge.dev/projectforge/app/file"
	"projectforge.dev/projectforge/app/lib/types"
	"projectforge.dev/projectforge/app/util"
)

func Model(m *model.Model, args *model.Args, addHeader bool) (*file.File, error) {
	g := golang.NewFile(m.Package, []string{"app", m.Package}, strings.ToLower(m.Camel()))
	for _, imp := range helper.ImportsForTypes("go", m.Columns.Types()...) {
		g.AddImport(imp)
	}
	for _, imp := range helper.ImportsForTypes("string", m.PKs().Types()...) {
		g.AddImport(imp)
	}
	if len(m.PKs()) > 1 {
		g.AddImport(helper.ImpFmt)
	}
	g.AddImport(helper.ImpAppUtil, helper.ImpSlices)
	for _, col := range m.Columns {
		if col.Type.Key() == types.KeyReference {
			ref, err := model.AsRef(col.Type)
			if err != nil {
				return nil, err
			}
			if ref.Pkg.Last() != m.Package {
				g.AddImport(golang.NewImport(golang.ImportTypeApp, ref.Pkg.ToPath()))
			}
		}
	}

	g.AddBlocks(modelStruct(m), modelConstructor(m), modelRandom(m))
	if b, err := modelFromMap(g, m); err == nil {
		g.AddBlocks(b)
	} else {
		return nil, err
	}
	g.AddBlocks(modelClone(m), modelString(m), modelWebPath(m), modelDiff(m, g), modelToData(m, m.Columns, ""))
	if m.IsRevision() {
		hc := m.HistoryColumns(false)
		g.AddBlocks(modelToData(m, hc.Const, "Core"), modelToData(m, hc.Var, hc.Col.Proper()))
	}
	g.AddBlocks(modelArray(m), modelArrayGet(m), modelArrayClone(m))
	return g.Render(addHeader)
}

func modelStruct(m *model.Model) *golang.Block {
	ret := golang.NewBlock(m.Proper(), "struct")
	ret.W("type %s struct {", m.Proper())
	maxColLength := util.StringArrayMaxLength(m.Columns.CamelNames())
	maxTypeLength := m.Columns.MaxGoKeyLength(m.Package)
	for _, c := range m.Columns {
		suffix := ""
		if c.Nullable || c.HasTag("omitempty") {
			suffix = ",omitempty"
		}
		goType := util.StringPad(c.ToGoType(m.Package), maxTypeLength)
		ret.W("\t%s %s `json:%q`", util.StringPad(c.Proper(), maxColLength), goType, c.Camel()+suffix)
	}
	ret.W("}")
	return ret
}

func modelConstructor(m *model.Model) *golang.Block {
	ret := golang.NewBlock("New"+m.Proper(), "func")
	ret.W("func New(%s) *%s {", m.PKs().Args(m.Package), m.Proper())
	ret.W("\treturn &%s{%s}", m.Proper(), m.PKs().Refs())
	ret.W("}")
	return ret
}

func modelDiff(m *model.Model, g *golang.File) *golang.Block {
	ret := golang.NewBlock("Diff"+m.Proper(), "func")
	ret.W("func (%s *%s) Diff(%sx *%s) util.Diffs {", m.FirstLetter(), m.Proper(), m.FirstLetter(), m.Proper())
	ret.W("\tvar diffs util.Diffs")
	for _, col := range m.Columns {
		if col.HasTag("updated") {
			continue
		}
		l := fmt.Sprintf("%s.%s", m.FirstLetter(), col.Proper())
		r := fmt.Sprintf("%sx.%s", m.FirstLetter(), col.Proper())
		switch col.Type.Key() {
		case types.KeyAny:
			ret.W("\tdiffs = append(diffs, util.DiffObjects(%s, %s, %q)...)", l, r, col.Camel())
		case types.KeyBool:
			g.AddImport(helper.ImpFmt)
			ret.W("\tif %s != %s {", l, r)
			ret.W("\t\tdiffs = append(diffs, util.NewDiff(%q, fmt.Sprint(%s), fmt.Sprint(%s)))", col.Camel(), l, r)
			ret.W("\t}")
		case types.KeyInt:
			g.AddImport(helper.ImpFmt)
			ret.W("\tif %s != %s {", l, r)
			ret.W("\t\tdiffs = append(diffs, util.NewDiff(%q, fmt.Sprint(%s), fmt.Sprint(%s)))", col.Camel(), l, r)
			ret.W("\t}")
		case types.KeyMap, types.KeyValueMap, types.KeyReference:
			ret.W("\tdiffs = append(diffs, util.DiffObjects(%s, %s, %q)...)", l, r, col.Camel())
		case types.KeyString:
			ret.W("\tif %s != %s {", l, r)
			ret.W("\t\tdiffs = append(diffs, util.NewDiff(%q, %s, %s))", col.Camel(), l, r)
			ret.W("\t}")
		case types.KeyTimestamp, types.KeyUUID:
			if col.Nullable {
				ret.W("\tif (%s == nil && %s != nil) || (%s != nil && %s == nil) || (%s != nil && %s != nil && *%s != *%s) {", l, r, l, r, l, r, l, r)
				g.AddImport(helper.ImpFmt)
				ret.W("\t\tdiffs = append(diffs, util.NewDiff(%q, fmt.Sprint(%s), fmt.Sprint(%s))) // nolint:gocritic // it's nullable", col.Camel(), l, r)
			} else {
				ret.W("\tif %s != %s {", l, r)
				ret.W("\t\tdiffs = append(diffs, util.NewDiff(%q, %s.String(), %s.String()))", col.Camel(), l, r)
			}
			ret.W("\t}")
		default:
			ret.W("\tTODO: %s", col.Type.Key())
		}
	}
	ret.W("\treturn diffs")
	ret.W("}")
	return ret
}

func modelToData(m *model.Model, cols model.Columns, suffix string) *golang.Block {
	ret := golang.NewBlock(m.Proper(), "func")
	ret.W("func (%s *%s) ToData%s() []any {", m.FirstLetter(), m.Proper(), suffix)
	refs := make([]string, 0, len(cols))
	for _, c := range cols {
		// switch c.Type.Key() {
		//   case types.KeyAny, types.KeyMap:
		//     ret.W("\t%sArg := util.ToJSONBytes(%s.%s, true)", c.Camel(), m.FirstLetter(), c.Proper())
		//     refs = append(refs, fmt.Sprintf("%sArg", c.Camel()))
		//   default:
		refs = append(refs, fmt.Sprintf("%s.%s", m.FirstLetter(), c.Proper()))
		// }
	}
	ret.W("\treturn []any{%s}", strings.Join(refs, ", "))
	ret.W("}")
	return ret
}
