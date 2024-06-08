package gomodel

import (
	"fmt"
	"strings"

	"github.com/pkg/errors"

	"projectforge.dev/projectforge/app/file"
	"projectforge.dev/projectforge/app/lib/metamodel/enum"
	"projectforge.dev/projectforge/app/lib/metamodel/model"
	"projectforge.dev/projectforge/app/lib/types"
	"projectforge.dev/projectforge/app/project/export/files/helper"
	"projectforge.dev/projectforge/app/project/export/golang"
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
	cols := m.Columns.NotDerived().WithoutTags("created", "updated")
	pks := cols.PKs()
	nonPKs := cols.NonPKs()

	ret := golang.NewBlock(m.Package+"FromForm", "func")
	ret.W("func FromMap(m util.ValueMap, setPK bool) (*%s, util.ValueMap, error) {", m.Proper())
	ret.W("\tret := &%s{}", m.Proper())
	ret.W("\textra := util.ValueMap{}")
	ret.W("\tfor k, v := range m {")
	ret.W("\t\tvar err error")
	ret.W("\t\tswitch k {")
	for _, col := range pks {
		ret.W("\t\tcase %q:", col.Camel())
		ret.W("\t\t\tif setPK {")
		if err := forCol(g, ret, 4, m, enums, col); err != nil {
			return nil, err
		}
		ret.W("\t\t\t}")
	}
	for _, col := range nonPKs {
		ret.W("\t\tcase %q:", col.Camel())
		if err := forCol(g, ret, 3, m, enums, col); err != nil {
			return nil, err
		}
	}
	ret.W("\t\tdefault:")
	ret.W("\t\t\textra[k] = v")
	ret.W("\t\t}")
	ret.W("\t\tif err != nil {")
	ret.W("\t\t\treturn nil, nil, err")
	ret.W("\t\t}")
	ret.W("\t}")
	ret.W("\t// $PF_SECTION_START(extrachecks)$")
	ret.W("\t// $PF_SECTION_END(extrachecks)$")
	ret.W("\treturn ret, extra, nil")
	ret.W("}")

	return ret, nil
}

func forCol(g *golang.File, ret *golang.Block, indent int, m *model.Model, enums enum.Enums, col *model.Column) error {
	ind := util.StringRepeat("\t", indent)
	catchErr := func(s string) {
		ret.W(ind + "if " + s + " != nil {")
		ret.W(ind + "\treturn nil, nil, " + s)
		ret.W(ind + "}")
	}
	parseCall := "ret%s, err := m.Parse%s(k, true, true)"
	parseMsg := "ret.%s, err = m.Parse%s(k, true, true)"
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
		ret.W(ind+parseCall, col.Proper(), col.ToGoMapParse())
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
			ret.W(ind+parseCall, col.Proper(), col.ToGoMapParse())
			catchErr("err")
			eRef := e.Proper()
			if e.PackageWithGroup("") != m.PackageWithGroup("") {
				eRef = e.Package + "." + eRef
			}
			ret.W(ind+"ret.%s = %sParse(nil, ret%s...)", col.Proper(), eRef, col.Proper())
		} else {
			ret.W(ind+parseMsg, col.Proper(), col.ToGoMapParse())
			catchErr("err")
		}
	case col.Nullable || col.Type.Scalar():
		ret.W(ind+parseMsg, col.Proper(), col.ToGoMapParse())
	default:
		ret.W(ind+"ret%s, e := m.Parse%s(k, true, true)", col.Proper(), col.ToGoMapParse())
		catchErr("e")
		ret.W(ind+"if ret%s != nil {", col.Proper())
		ret.W(ind+"\tret.%s = *ret%s", col.Proper(), col.Proper())
		ret.W(ind + "}")
	}
	return nil
}
