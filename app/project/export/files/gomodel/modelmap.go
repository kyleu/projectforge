package gomodel

import (
	"fmt"
	"strings"

	"github.com/pkg/errors"
	"github.com/samber/lo"

	"projectforge.dev/projectforge/app/file"
	"projectforge.dev/projectforge/app/lib/metamodel/enum"
	"projectforge.dev/projectforge/app/lib/metamodel/model"
	"projectforge.dev/projectforge/app/lib/types"
	"projectforge.dev/projectforge/app/project/export/files/helper"
	"projectforge.dev/projectforge/app/project/export/golang"
	"projectforge.dev/projectforge/app/util"
)

func ModelMap(m *model.Model, args *model.Args, linebreak string) (*file.File, error) {
	g := golang.NewFile(m.Package, []string{"app", m.PackageWithGroup("")}, strings.ToLower(m.Camel())+"map")
	g.AddImport(helper.ImpAppUtil)
	imps, err := helper.SpecialImports(m.Columns, m.PackageWithGroup(""), args.Models, args.Enums)
	if err != nil {
		return nil, err
	}
	g.AddImport(imps...)
	imps, err = helper.EnumImports(m.Columns.Types(), m.PackageWithGroup(""), args.Models, args.Enums)
	if err != nil {
		return nil, err
	}
	g.AddImport(imps...)
	g.AddImport(m.Imports.Supporting("map")...)
	g.AddBlocks(modelToMap(g, m, args.Enums, args.Database))
	if b, e := modelFromMap(g, m, args.Models, args.Enums, args.Database); e == nil {
		g.AddBlocks(b)
	} else {
		return nil, e
	}
	return g.Render(linebreak)
}

func modelToMap(g *golang.File, m *model.Model, enums enum.Enums, database string) *golang.Block {
	ret := golang.NewBlock(m.Package+"ToMap", "func")
	ret.WF("func (%s *%s) ToMap() util.ValueMap {", m.FirstLetter(), m.Proper())
	content := strings.Join(lo.Map(m.Columns, func(col *model.Column, _ int) string {
		return fmt.Sprintf(`%q: %s.%s`, col.Camel(), m.FirstLetter(), col.Proper())
	}), ", ")
	ret.W("\treturn util.ValueMap{" + content + "}")
	ret.W("}")

	_ = ` {
		return util.ValueMap{"id": m.ID, "t": m.T, "data": m.Data, "occurred": m.Occurred}
	}`
	return ret
}

func modelFromMap(g *golang.File, m *model.Model, models model.Models, enums enum.Enums, database string) (*golang.Block, error) {
	cols := m.Columns.NotDerived().WithoutTags("created", "updated")
	pks := cols.PKs()
	nonPKs := cols.NonPKs()

	ret := golang.NewBlock(m.Package+"FromMap", "func")
	ret.WF("func %sFromMap(m util.ValueMap, setPK bool) (*%s, util.ValueMap, error) {", m.Proper(), m.Proper())
	ret.WF("\tret := &%s{}", m.Proper())
	ret.W("\textra := util.ValueMap{}")
	ret.W("\tfor k, v := range m {")
	ret.W("\t\tvar err error")
	ret.W("\t\tswitch k {")
	for _, col := range pks {
		ret.WF("\t\tcase %q:", col.CamelNoReplace())
		ret.W("\t\t\tif setPK {")
		if err := forCol(g, ret, 4, m, models, enums, col); err != nil {
			return nil, err
		}
		ret.W("\t\t\t}")
	}
	for _, col := range nonPKs {
		ret.WF("\t\tcase %q:", col.CamelNoReplace())
		if err := forCol(g, ret, 3, m, models, enums, col); err != nil {
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

//nolint:gocognit
func forCol(g *golang.File, ret *golang.Block, indent int, m *model.Model, models model.Models, enums enum.Enums, col *model.Column) error {
	ind := util.StringRepeat("\t", indent)
	catchErr := func(s string) {
		ret.W(ind + "if " + s + " != nil {")
		ret.W(ind + "\treturn nil, nil, " + s)
		ret.W(ind + "}")
	}
	colMP := func(msg string) error {
		mp, err := col.ToGoMapParse()
		if err != nil {
			return err
		}
		ret.WF(msg, col.Proper(), mp)
		return nil
	}
	parseCall := "ret%s, err := m.Parse%s(k, true, true)"
	parseMsg := "ret.%s, err = m.Parse%s(k, true, true)"
	switch {
	case col.Type.Key() == types.KeyAny:
		ret.WF(ind+"ret.%s = m[%q]", col.Proper(), col.Camel())
	case col.Type.Key() == types.KeyReference:
		ret.WF(ind+"tmp%s, err := m.ParseString(%q, true, true)", col.Proper(), col.Camel())
		catchErr("err")
		ref, err := helper.LoadRef(col, models)
		if err != nil {
			return errors.Wrap(err, "invalid ref")
		}
		ret.WF(ind+"%sArg := %s{}", col.Camel(), ref.LastAddr(ref.Pkg.Last() != m.Package))
		if ref.Pkg.Last() != m.Package {
			g.AddImport(model.NewImport(model.ImportTypeApp, ref.Pkg.ToPath()))
		}
		ret.WF(ind+"err = util.FromJSON([]byte(tmp%s), %sArg)", col.Proper(), col.Camel())
		catchErr("err")
		ret.WF(ind+"ret.%s = %sArg", col.Proper(), col.Camel())
	case col.Type.Key() == types.KeyEnum:
		e, err := model.AsEnumInstance(col.Type, enums)
		if err != nil {
			return err
		}
		if err := colMP(ind + parseCall); err != nil {
			return err
		}
		catchErr("err")
		var enumRef string
		if e.Simple() {
			enumRef = fmt.Sprintf("%s(ret%s)", e.Proper(), col.Proper())
		} else {
			enumRef = fmt.Sprintf("All%s.Get(ret%s, nil)", e.ProperPlural(), col.Proper())
		}
		if e.PackageWithGroup("") == m.PackageWithGroup("") {
			ret.WF(ind+"ret.%s = %s", col.Proper(), enumRef)
		} else {
			ret.WF(ind+"ret.%s = %s.%s", col.Proper(), e.Package, enumRef)
		}
	case col.Type.Key() == types.KeyJSON:
		if err := colMP(ind + "ret.%s, err = m.Parse%s(k, true, true)"); err != nil {
			return err
		}
	case col.Type.Key() == types.KeyList:
		if e, _ := model.AsEnumInstance(col.Type.ListType(), enums); e != nil {
			if err := colMP(ind + parseCall); err != nil {
				return err
			}
			catchErr("err")
			eRef := e.Proper()
			if e.PackageWithGroup("") != m.PackageWithGroup("") {
				eRef = e.Package + "." + eRef
			}
			ret.WF(ind+"ret.%s = %sParse(nil, ret%s...)", col.Proper(), eRef, col.Proper())
		} else {
			if err := colMP(ind + parseMsg); err != nil {
				return err
			}
			catchErr("err")
		}
	case col.Nullable || col.Type.Scalar():
		if err := colMP(ind + parseMsg); err != nil {
			return err
		}
	default:
		if err := colMP(ind + "ret%s, e := m.Parse%s(k, true, true)"); err != nil {
			return err
		}
		catchErr("e")
		ret.WF(ind+"if ret%s != nil {", col.Proper())
		ret.WF(ind+"\tret.%s = *ret%s", col.Proper(), col.Proper())
		ret.W(ind + "}")
	}
	return nil
}
