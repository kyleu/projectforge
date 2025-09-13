package goevent

import (
	"strings"

	"projectforge.dev/projectforge/app/file"
	"projectforge.dev/projectforge/app/lib/metamodel"
	"projectforge.dev/projectforge/app/lib/metamodel/model"
	"projectforge.dev/projectforge/app/project/export/files/gohelper"
	"projectforge.dev/projectforge/app/project/export/files/helper"
	"projectforge.dev/projectforge/app/project/export/golang"
)

func EventDiff(evt *model.Event, args *metamodel.Args, linebreak string) (*file.File, error) {
	g := golang.NewFile(evt.Package, []string{"app", evt.PackageWithGroup("")}, strings.ToLower(evt.Camel())+"diff")
	g.AddImport(helper.ImpAppUtil)
	g.AddImport(evt.Imports.Supporting("diff")...)

	mdiff, err := gohelper.DiffBlock(g, evt.Columns, evt, args.Enums)
	if err != nil {
		return nil, err
	}
	g.AddBlocks(mdiff)

	return g.Render(linebreak)
}
