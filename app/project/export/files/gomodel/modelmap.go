package gomodel

import (
	"strings"

	"projectforge.dev/projectforge/app/file"
	"projectforge.dev/projectforge/app/lib/metamodel"
	"projectforge.dev/projectforge/app/lib/metamodel/model"
	"projectforge.dev/projectforge/app/project/export/files/gohelper"
	"projectforge.dev/projectforge/app/project/export/files/helper"
	"projectforge.dev/projectforge/app/project/export/golang"
)

func ModelMap(m *model.Model, args *metamodel.Args, linebreak string) (*file.File, error) {
	g := golang.NewFile(m.Package, []string{"app", m.PackageWithGroup("")}, strings.ToLower(m.Camel())+"map")
	g.AddImport(helper.ImpAppUtil)
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
	g.AddImport(m.Imports.Supporting("map")...)
	g.AddBlocks(gohelper.ToMap(m, m.Columns))
	if b, e := gohelper.FromMap(g, m, m.Columns, args); e == nil {
		g.AddBlocks(b)
	} else {
		return nil, e
	}
	g.AddBlocks(gohelper.ToOrderedMap(m, m.Columns.NotDerived()))
	return g.Render(linebreak)
}
