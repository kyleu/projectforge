package goevent

import (
	"strings"

	"projectforge.dev/projectforge/app/file"
	"projectforge.dev/projectforge/app/lib/metamodel"
	"projectforge.dev/projectforge/app/lib/metamodel/model"
	"projectforge.dev/projectforge/app/project/export/files/gohelper"
	"projectforge.dev/projectforge/app/project/export/golang"
)

func Events(m *model.Event, args *metamodel.Args, goVersion string, linebreak string) (*file.File, error) {
	name := strings.ToLower(m.CamelPlural())
	if name == strings.ToLower(m.Camel()) {
		name += "_array"
	}
	g := golang.NewFile(m.PackageName(), []string{"app", m.PackageWithGroup("")}, name)
	g.AddImport(m.Imports.Supporting("array")...)
	gohelper.Array(g, m, m.Columns, nil, args, goVersion, linebreak)
	return g.Render(linebreak)
}
