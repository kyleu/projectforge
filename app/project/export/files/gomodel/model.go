package gomodel

import (
	"strings"

	"github.com/samber/lo"

	"projectforge.dev/projectforge/app/file"
	"projectforge.dev/projectforge/app/lib/metamodel"
	"projectforge.dev/projectforge/app/lib/metamodel/enum"
	"projectforge.dev/projectforge/app/lib/metamodel/model"
	"projectforge.dev/projectforge/app/lib/types"
	"projectforge.dev/projectforge/app/project/export/files/gohelper"
	"projectforge.dev/projectforge/app/project/export/files/helper"
	"projectforge.dev/projectforge/app/project/export/golang"
	"projectforge.dev/projectforge/app/util"
)

func Model(m *model.Model, args *metamodel.Args, linebreak string) (*file.File, error) {
	g := golang.NewFile(m.Package, []string{"app", m.PackageWithGroup("")}, strings.ToLower(m.Camel()))
	lo.ForEach(helper.ImportsForTypes("go", "", m.Columns.Types()...), func(imp *model.Import, _ int) {
		g.AddImport(imp)
	})
	lo.ForEach(helper.ImportsForTypes(types.KeyString, "", m.PKs().Types()...), func(imp *model.Import, _ int) {
		g.AddImport(imp)
	})
	g.AddImport(helper.ImpAppUtil, helper.ImpAppSvc)
	imps, err := helper.SpecialImports(m.Columns, m.PackageWithGroup(""), args)
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

	if !m.SkipController() {
		g.AddBlocks(defaultRoute(m), routeMethod())
	}
	g.AddBlocks(typeAssert(m))

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
	str, err := modelStruct(m, args)
	if err != nil {
		return nil, err
	}
	c, err := modelConstructor(m, args)
	if err != nil {
		return nil, err
	}
	g.AddBlocks(str, c)

	g.AddBlocks(gohelper.BlockClone(g, m.Columns.NotDerived(), m), gohelper.BlockString(g, m.Columns, m), gohelper.BlockTitle(g, m.Columns, m))
	if len(m.PKs()) > 1 {
		if pk, e := modelToPK(m, args.Enums); e == nil {
			g.AddBlocks(pk)
		} else {
			return nil, err
		}
	}

	rnd, err := gohelper.BlockRandom(m.Columns, m, args.Enums)
	if err != nil {
		return nil, err
	}
	g.AddBlocks(rnd, gohelper.BlockStrings(g, m.Columns, m), gohelper.BlockToCSV(m))

	fd, err := gohelper.BlockFieldDescs(m.Columns, m)
	if err != nil {
		return nil, err
	}
	if !m.SkipController() {
		g.AddBlocks(modelWebPath(g, m), modelBreadcrumb(m))
	}
	g.AddBlocks(gohelper.BlockToData(m, m.Columns.NotDerived(), "", args.Database), fd)
	return g.Render(linebreak)
}

func typeAssert(m *model.Model) *golang.Block {
	ret := golang.NewBlock("assert", "type")
	ret.WF("var _ svc.Model = (*%s)(nil)", m.Proper())
	return ret
}

func defaultRoute(m *model.Model) *golang.Block {
	ret := golang.NewBlock("DefaultRoute", "const")
	ret.WF("const DefaultRoute = %q", "/"+m.Route())
	return ret
}

func routeMethod() *golang.Block {
	ret := golang.NewBlock("Route", "func")
	ret.W("func Route(paths ...string) string {")
	ret.W("\tif len(paths) == 0 {")
	ret.W("\t\tpaths = []string{DefaultRoute}")
	ret.W("\t}")
	ret.W("\treturn util.StringPath(paths...)")
	ret.W("}")
	return ret
}

func modelToPK(m *model.Model, _ enum.Enums) (*golang.Block, error) {
	ret := golang.NewBlock("PK", "struct")
	ret.WF("func (%s *%s) ToPK() *PK {", m.FirstLetter(), m.Proper())
	ret.W("\treturn &PK{")
	pks := m.PKs()
	maxColLength := pks.MaxCamelLength() + 1
	for _, c := range pks {
		ret.WF("\t\t%s %s.%s,", util.StringPad(c.Proper()+":", maxColLength), m.FirstLetter(), c.Proper())
	}
	ret.W("\t}")
	ret.W("}")
	return ret, nil
}

func ModelDiff(m *model.Model, args *metamodel.Args, linebreak string) (*file.File, error) {
	g := golang.NewFile(m.Package, []string{"app", m.PackageWithGroup("")}, strings.ToLower(m.Camel())+"diff")
	g.AddImport(helper.ImpAppUtil)
	g.AddImport(m.Imports.Supporting("diff")...)
	mdiff, err := gohelper.DiffBlock(g, m.Columns, m, args.Enums)
	if err != nil {
		return nil, err
	}
	g.AddBlocks(mdiff)
	return g.Render(linebreak)
}
