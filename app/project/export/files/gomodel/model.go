package gomodel

import (
	"fmt"
	"strings"

	"projectforge.dev/projectforge/app/project/export/enum"

	"projectforge.dev/projectforge/app/file"
	"projectforge.dev/projectforge/app/lib/types"
	"projectforge.dev/projectforge/app/project/export/files/helper"
	"projectforge.dev/projectforge/app/project/export/golang"
	"projectforge.dev/projectforge/app/project/export/model"
	"projectforge.dev/projectforge/app/util"
)

func Model(m *model.Model, args *model.Args, addHeader bool) (*file.File, error) {
	g := golang.NewFile(m.Package, []string{"app", m.PackageWithGroup("")}, strings.ToLower(m.Camel()))
	for _, imp := range helper.ImportsForTypes("go", args.Enums, m.Columns.Types()...) {
		g.AddImport(imp)
	}
	for _, imp := range helper.ImportsForTypes("string", args.Enums, m.PKs().Types()...) {
		g.AddImport(imp)
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

	if len(m.PKs()) > 1 {
		pk, err := modelPK(m, args.Enums)
		if err != nil {
			return nil, err
		}
		g.AddBlocks(pk)
	}
	str, err := modelStruct(m, args.Enums)
	if err != nil {
		return nil, err
	}
	c, err := modelConstructor(m, args.Enums)
	if err != nil {
		return nil, err
	}
	g.AddBlocks(str, c, modelRandom(m, args.Enums))
	if b, e := modelFromMap(g, m, args.Enums); e == nil {
		g.AddBlocks(b)
	} else {
		return nil, err
	}
	g.AddBlocks(modelClone(m), modelString(g, m), modelTitle(m), modelWebPath(g, m), modelDiff(m, g), modelToData(m, m.Columns, ""))
	if m.IsRevision() {
		hc := m.HistoryColumns(false)
		g.AddBlocks(modelToData(m, hc.Const, "Core"), modelToData(m, hc.Var, hc.Col.Proper()))
	}
	ag, err := modelArrayGet(g, m, args.Enums)
	if err != nil {
		return nil, err
	}
	g.AddBlocks(modelArray(m), ag, modelArrayClone(m))
	return g.Render(addHeader)
}

func modelPK(m *model.Model, enums enum.Enums) (*golang.Block, error) {
	ret := golang.NewBlock("PK", "struct")
	ret.W("type PK struct {")
	pks := m.PKs()
	maxColLength := pks.MaxCamelLength()
	maxTypeLength := pks.MaxGoTypeLength(m.Package, enums)
	for _, c := range pks {
		gt, err := c.ToGoType(m.Package, enums)
		if err != nil {
			return nil, err
		}
		goType := util.StringPad(gt, maxTypeLength)
		ret.W("\t%s %s `json:%q`", util.StringPad(c.Proper(), maxColLength), goType, c.Camel()+modelJSONSuffix(c))
	}
	ret.W("}")
	return ret, nil
}

func modelStruct(m *model.Model, enums enum.Enums) (*golang.Block, error) {
	ret := golang.NewBlock(m.Proper(), "struct")
	ret.W("type %s struct {", m.Proper())
	maxColLength := m.Columns.MaxCamelLength()
	maxTypeLength := m.Columns.MaxGoTypeLength(m.Package, enums)
	for _, c := range m.Columns {
		gt, err := c.ToGoType(m.Package, enums)
		if err != nil {
			return nil, err
		}
		goType := util.StringPad(gt, maxTypeLength)
		ret.W("\t%s %s `json:%q`", util.StringPad(c.Proper(), maxColLength), goType, c.Camel()+modelJSONSuffix(c))
	}
	ret.W("}")
	return ret, nil
}

func modelJSONSuffix(c *model.Column) string {
	if c.Nullable || c.HasTag("omitempty") {
		return ",omitempty"
	}
	return ""
}

func modelConstructor(m *model.Model, enums enum.Enums) (*golang.Block, error) {
	ret := golang.NewBlock("New"+m.Proper(), "func")
	args, err := m.PKs().Args(m.Package, enums)
	if err != nil {
		return nil, err
	}
	ret.W("func New(%s) *%s {", args, m.Proper())
	ret.W("\treturn &%s{%s}", m.Proper(), m.PKs().Refs())
	ret.W("}")
	return ret, nil
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
