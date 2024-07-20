package gomodel

import (
	"fmt"
	"strings"

	"github.com/samber/lo"

	"projectforge.dev/projectforge/app/file"
	"projectforge.dev/projectforge/app/lib/metamodel/enum"
	"projectforge.dev/projectforge/app/lib/metamodel/model"
	"projectforge.dev/projectforge/app/lib/types"
	"projectforge.dev/projectforge/app/project/export/files/helper"
	"projectforge.dev/projectforge/app/project/export/golang"
	"projectforge.dev/projectforge/app/util"
)

func Model(m *model.Model, args *model.Args, linebreak string) (*file.File, error) {
	g := golang.NewFile(m.Package, []string{"app", m.PackageWithGroup("")}, strings.ToLower(m.Camel()))
	lo.ForEach(helper.ImportsForTypes("go", "", m.Columns.Types()...), func(imp *model.Import, _ int) {
		g.AddImport(imp)
	})
	lo.ForEach(helper.ImportsForTypes(types.KeyString, "", m.PKs().Types()...), func(imp *model.Import, _ int) {
		g.AddImport(imp)
	})
	g.AddImport(helper.ImpAppUtil, helper.ImpAppSvc, helper.ImpPath)
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
	g.AddImport(m.Imports.Supporting("model")...)

	g.AddBlocks(defaultRoute(m), routeMethod(m), typeAssert(m))

	if len(m.PKs()) > 1 {
		pk, e := modelPK(m, args.Enums)
		if e != nil {
			return nil, e
		}
		pkString, e := modelPKString(m)
		if e != nil {
			return nil, e
		}
		g.AddBlocks(pk, pkString)
	}
	str, err := modelStruct(m, args.Enums)
	if err != nil {
		return nil, err
	}
	c, err := modelConstructor(m, args.Enums)
	if err != nil {
		return nil, err
	}
	g.AddBlocks(str, c)

	g.AddBlocks(modelClone(m), modelString(g, m), modelTitle(g, m))
	if len(m.PKs()) > 1 {
		if pk, e := modelToPK(m, args.Enums); e == nil {
			g.AddBlocks(pk)
		} else {
			return nil, err
		}
	}

	rnd, err := modelRandom(m, args.Enums)
	if err != nil {
		return nil, err
	}
	g.AddBlocks(rnd, modelStrings(g, m), modelToCSV(m))

	fd, err := modelFieldDescs(m)
	if err != nil {
		return nil, err
	}

	g.AddBlocks(modelWebPath(g, m), modelToData(m, m.Columns.NotDerived(), "", args.Database), fd)
	return g.Render(linebreak)
}

func typeAssert(m *model.Model) *golang.Block {
	ret := golang.NewBlock("assert", "type")
	ret.W("var _ svc.Model = (*%s)(nil)", m.Proper())
	return ret
}

func defaultRoute(m *model.Model) *golang.Block {
	ret := golang.NewBlock("DefaultRoute", "const")
	ret.W("const DefaultRoute = %q", "/"+m.Route())
	return ret
}

func routeMethod(m *model.Model) *golang.Block {
	ret := golang.NewBlock("Route", "func")
	ret.W("func Route(paths ...string) string {")
	ret.W("\tif len(paths) == 0 {")
	ret.W("\t\tpaths = []string{DefaultRoute}")
	ret.W("\t}")
	ret.W("\treturn path.Join(paths...)")
	ret.W("}")
	return ret
}

func modelToPK(m *model.Model, _ enum.Enums) (*golang.Block, error) {
	ret := golang.NewBlock("PK", "struct")
	ret.W("func (%s *%s) ToPK() *PK {", m.FirstLetter(), m.Proper())
	ret.W("\treturn &PK{")
	pks := m.PKs()
	maxColLength := pks.MaxCamelLength() + 1
	for _, c := range pks {
		ret.W("\t\t%s %s.%s,", util.StringPad(c.Proper()+":", maxColLength), m.FirstLetter(), c.Proper())
	}
	ret.W("\t}")
	ret.W("}")
	return ret, nil
}

func modelStruct(m *model.Model, enums enum.Enums) (*golang.Block, error) {
	ret := golang.NewBlock(m.Proper(), "struct")
	ret.W("type %s struct {", m.Proper())
	cols := m.Columns.NotDerived()
	maxColLength := cols.MaxCamelLength()
	maxTypeLength := cols.MaxGoTypeLength(m.Package, enums)
	for _, c := range cols {
		gt, err := c.ToGoType(m.Package, enums)
		if err != nil {
			return nil, err
		}
		goType := util.StringPad(gt, maxTypeLength)
		var tag string
		if c.JSON == "" {
			tag = fmt.Sprintf("json:%q", c.CamelNoReplace()+modelJSONSuffix(c))
		} else {
			tag = fmt.Sprintf("json:%q", c.JSON+modelJSONSuffix(c))
		}
		if c.Validation != "" {
			tag += fmt.Sprintf(",validate:%q", c.Validation)
		}
		if c.Example != "" {
			tag += fmt.Sprintf(",fake:%q", c.Example)
		}
		ret.W("\t%s %s `%s`", util.StringPad(c.Proper(), maxColLength), goType, tag)
	}
	ret.W("}")
	return ret, nil
}

func modelJSONSuffix(_ *model.Column) string {
	return ",omitempty"
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

func modelToData(m *model.Model, cols model.Columns, suffix string, database string) *golang.Block {
	ret := golang.NewBlock(m.Proper(), "func")
	ret.W("func (%s *%s) ToData%s() []any {", m.FirstLetter(), m.Proper(), suffix)
	calls := lo.Map(cols, func(c *model.Column, _ int) string {
		tk := c.Type.Key()
		switch {
		case (tk == types.KeyAny || tk == types.KeyList || tk == types.KeyMap || tk == types.KeyReference) && (helper.SimpleJSON(database)):
			return fmt.Sprintf("util.ToJSON(%s.%s),", m.FirstLetter(), c.Proper())
		default:
			return fmt.Sprintf("%s.%s,", m.FirstLetter(), c.Proper())
		}
	})
	lines := util.JoinLines(calls, " ", 120)
	if len(lines) == 1 && len(lines[0]) < 100 {
		ret.W("\treturn []any{%s}", strings.TrimSuffix(lines[0], ","))
	} else {
		ret.W("\treturn []any{")
		lo.ForEach(lines, func(l string, _ int) {
			ret.W("\t\t%s", l)
		})
		ret.W("\t}")
	}
	ret.W("}")
	return ret
}
