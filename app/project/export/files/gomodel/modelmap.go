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

func ModelMap(m *model.Model, args *model.Args, addHeader bool, linebreak string) (*file.File, error) {
	g := golang.NewFile(m.Package, []string{"app", m.PackageWithGroup("")}, strings.ToLower(m.Camel())+"map")
	g.AddImport(helper.ImpAppUtil)
	imps, err := helper.SpecialImports(m.Columns, m.PackageWithGroup(""), args.Enums)
	if err != nil {
		return nil, err
	}
	g.AddImport(imps...)
	imps, err = helper.EnumImports(m.Columns.Types(), m.PackageWithGroup(""), args.Enums)
	if err != nil {
		return nil, err
	}
	g.AddImport(imps...)
	if b, e := modelFromMap(g, m, args.Enums, args.Database); e == nil {
		g.AddBlocks(b)
	} else {
		return nil, err
	}
	return g.Render(addHeader, linebreak)
}

func modelFromMap(g *golang.File, m *model.Model, enums enum.Enums, database string) (*golang.Block, error) {
	ret := golang.NewBlock(m.Package+"FromForm", "func")
	ret.W("func FromMap(m util.ValueMap, setPK bool) (*%s, error) {", m.Proper())
	ret.W("\tret := &%s{}", m.Proper())
	cols := m.Columns.WithoutTags("created", "updated")
	needsErr := lo.ContainsBy(cols, func(c *model.Column) bool {
		return c.NeedsErr(m.Name, database)
	})
	if needsErr {
		ret.W("\tvar err error")
	}
	ret.W("\tif setPK {")
	err := forCols(g, ret, 2, m, enums, cols.PKs()...)
	if err != nil {
		return nil, err
	}
	ret.W("\t\t// $PF_SECTION_START(pkchecks)$")
	ret.W("\t\t// $PF_SECTION_END(pkchecks)$")
	ret.W("\t}")
	err = forCols(g, ret, 1, m, enums, cols.NonPKs()...)
	if err != nil {
		return nil, err
	}
	ret.W("\t// $PF_SECTION_START(extrachecks)$")
	ret.W("\t// $PF_SECTION_END(extrachecks)$")
	ret.W("\treturn ret, nil")
	ret.W("}")

	return ret, nil
}

func forCols(g *golang.File, ret *golang.Block, indent int, m *model.Model, enums enum.Enums, cols ...*model.Column) error {
	ind := util.StringRepeat("\t", indent)
	catchErr := func(s string) {
		ret.W(ind + "if " + s + " != nil {")
		ret.W(ind + "\treturn nil, " + s)
		ret.W(ind + "}")
	}
	parseCall := "ret%s, err := m.Parse%s(%q, true, true)"
	parseMsg := "ret.%s, err = m.Parse%s(%q, true, true)"
	for _, col := range cols {
		switch {
		case col.Type.Key() == types.KeyAny:
			ret.W(ind+"ret.%s = m[%q]", col.Proper(), col.Camel())
		case col.Type.Key() == types.KeyReference:
			ret.W(ind+"tmp%s, err := m.ParseString(%q, true, true)", col.Proper(), col.Camel())
			catchErr("err")

			ref, err := model.AsRef(col.Type)
			if err != nil {
				return errors.Wrap(err, "invalid ref")
			}
			if ref.Pkg.Last() == m.Package {
				ret.W("\t%sArg := &%s{}", col.Camel(), ref.K)
			} else {
				ret.W("\t%sArg := &%s.%s{}", col.Camel(), ref.Pkg.Last(), ref.K)
				g.AddImport(golang.NewImport(golang.ImportTypeApp, ref.Pkg.ToPath()))
			}
			ret.W(ind+"err = util.FromJSON([]byte(tmp%s), %sArg)", col.Proper(), col.Camel())
			catchErr("err")
			ret.W(ind+"ret.%s = %sArg", col.Proper(), col.Camel())
		case col.Type.Key() == types.KeyEnum:
			e, err := model.AsEnumInstance(col.Type, enums)
			if err != nil {
				return err
			}
			ret.W(ind+parseCall, col.Proper(), col.ToGoMapParse(), col.Camel())
			catchErr("err")
			var enumRef string
			if e.Simple() {
				enumRef = fmt.Sprintf("%s(ret%s)", e.Proper(), col.Proper())
			} else {
				enumRef = fmt.Sprintf("All%s.Get(ret%s, nil)", e.ProperPlural(), col.Proper())
			}
			if e.PackageWithGroup("") == m.PackageWithGroup("") {
				ret.W(ind+"ret.%s = %s", col.Proper(), enumRef)
			} else {
				ret.W(ind+"ret.%s = %s.%s", col.Proper(), e.Package, enumRef)
			}
		case col.Type.Key() == types.KeyList:
			lt := types.TypeAs[*types.List](col.Type)
			if e, _ := model.AsEnumInstance(lt.V, enums); e != nil {
				ret.W(ind+parseCall, col.Proper(), col.ToGoMapParse(), col.Camel())
				catchErr("err")
				eRef := e.Proper()
				if e.PackageWithGroup("") != m.PackageWithGroup("") {
					eRef = e.Package + "." + eRef
				}
				ret.W(ind+"ret.%s = %sParse(nil, ret%s...)", col.Proper(), eRef, col.Proper())
			} else {
				ret.W(ind+parseMsg, col.Proper(), col.ToGoMapParse(), col.Camel())
				catchErr("err")
			}
		case col.Nullable || col.Type.Scalar():
			ret.W(ind+parseMsg, col.Proper(), col.ToGoMapParse(), col.Camel())
			catchErr("err")
		default:
			ret.W(ind+"ret%s, e := m.Parse%s(%q, true, true)", col.Proper(), col.ToGoMapParse(), col.Camel())
			catchErr("e")
			ret.W(ind+"if ret%s != nil {", col.Proper())
			ret.W(ind+"\tret.%s = *ret%s", col.Proper(), col.Proper())
			ret.W(ind + "}")
		}
	}
	return nil
}
