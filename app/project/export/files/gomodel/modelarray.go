package gomodel

import (
	"strings"

	"github.com/samber/lo"

	"projectforge.dev/projectforge/app/file"
	"projectforge.dev/projectforge/app/lib/metamodel"
	"projectforge.dev/projectforge/app/lib/metamodel/model"
	"projectforge.dev/projectforge/app/lib/types"
	"projectforge.dev/projectforge/app/project/export/files/gohelper"
	"projectforge.dev/projectforge/app/project/export/files/helper"
	"projectforge.dev/projectforge/app/project/export/golang"
)

func Models(m *model.Model, args *metamodel.Args, goVersion string, linebreak string) (*file.File, error) {
	name := strings.ToLower(m.CamelPlural())
	if name == strings.ToLower(m.Camel()) {
		name += "_array"
	}
	g := golang.NewFile(m.PackageName(), []string{"app", m.PackageWithGroup("")}, name)
	lo.ForEach(helper.ImportsForTypes("go", "", m.IndexedColumns(true).Types()...), func(imp *model.Import, _ int) {
		g.AddImport(imp)
	})
	lo.ForEach(helper.ImportsForTypes(types.KeyString, "", m.PKs().Types()...), func(imp *model.Import, _ int) {
		g.AddImport(imp)
	})
	g.AddImport(m.Imports.Supporting("array")...)
	gohelper.Array(g, m, m.Columns, m.IndexedColumns(true), args, goVersion, linebreak)
	return g.Render(linebreak)
}
