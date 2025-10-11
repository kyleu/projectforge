package goevent

import (
	"projectforge.dev/projectforge/app/lib/metamodel"
	"projectforge.dev/projectforge/app/lib/metamodel/model"
	"projectforge.dev/projectforge/app/project/export/files/gohelper"
	"projectforge.dev/projectforge/app/project/export/files/helper"
	"projectforge.dev/projectforge/app/project/export/golang"
)

func eventMap(g *golang.File, evt *model.Event, args *metamodel.Args, linebreak string) error {
	g.AddImport(helper.ImpAppUtil)
	imps, err := helper.SpecialImports(evt.Columns, evt.PackageWithGroup(""), args)
	if err != nil {
		return err
	}
	g.AddImport(imps...)
	imps, err = helper.EnumImports(evt.Columns.Types(), evt.PackageWithGroup(""), args.Enums)
	if err != nil {
		return err
	}
	g.AddImport(imps...)
	g.AddImport(evt.Imports.Supporting("map")...)
	if b, e := gohelper.FromMap(g, evt, evt.Columns, args); e == nil {
		g.AddBlocks(b)
	} else {
		return e
	}
	g.AddBlocks(gohelper.ToMap(evt, evt.Columns))
	b := gohelper.ToOrderedMap(evt, evt.Columns.NotDerived())
	g.AddBlocks(b)
	return nil
}
