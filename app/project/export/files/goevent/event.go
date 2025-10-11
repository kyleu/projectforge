package goevent

import (
	"fmt"
	"strings"

	"github.com/samber/lo"

	"projectforge.dev/projectforge/app/file"
	"projectforge.dev/projectforge/app/lib/metamodel"
	"projectforge.dev/projectforge/app/lib/metamodel/model"
	"projectforge.dev/projectforge/app/lib/types"
	"projectforge.dev/projectforge/app/project/export/files/gohelper"
	"projectforge.dev/projectforge/app/project/export/files/helper"
	"projectforge.dev/projectforge/app/project/export/golang"
	"projectforge.dev/projectforge/app/util"
)

func Event(evt *model.Event, args *metamodel.Args, goVersion string, linebreak string) (*file.File, error) {
	g := golang.NewFile(evt.Package, []string{"app", evt.PackageWithGroup("")}, strings.ToLower(evt.Camel()))
	lo.ForEach(helper.ImportsForTypes("go", "", evt.Columns.Types()...), func(imp *model.Import, _ int) {
		g.AddImport(imp)
	})
	g.AddImport(helper.ImpAppUtil, helper.ImpLo)
	imps, err := helper.SpecialImports(evt.Columns, evt.PackageWithGroup(""), args)
	if err != nil {
		return nil, err
	}
	g.AddImport(imps...)
	imps, err = helper.EnumImports(evt.Columns.Types(), evt.PackageWithGroup(""), args.Enums)
	if err != nil {
		return nil, err
	}
	g.AddImport(imps...)
	g.AddImport(evt.Imports.Supporting("event")...)

	str, err := eventStruct(evt, args)
	if err != nil {
		return nil, err
	}
	g.AddBlocks(str)

	cn, err := eventConstructor(evt, args)
	if err != nil {
		return nil, err
	}
	g.AddBlocks(cn)

	g.AddBlocks(gohelper.BlockClone(g, evt.Columns, evt))
	if len(evt.Columns.PKs()) > 0 {
		g.AddBlocks(gohelper.BlockString(g, evt.Columns, evt), gohelper.BlockTitle(g, evt.Columns, evt))
	}
	err = eventMap(g, evt, args, linebreak)
	if err != nil {
		return nil, err
	}

	db, err := eventDiff(g, evt, args, linebreak)
	if err != nil {
		return nil, err
	}
	g.AddBlocks(db)

	g.AddBlocks(gohelper.BlockStrings(g, evt.Columns, evt), gohelper.BlockToCSV(evt))

	rnd, err := gohelper.BlockRandom(evt.Columns, evt, args.Enums)
	if err != nil {
		return nil, err
	}
	fd, err := gohelper.BlockFieldDescs(evt.Columns, evt)
	if err != nil {
		return nil, err
	}

	g.AddBlocks(eventToData(evt, evt.Columns.NotDerived(), "", args.Database), rnd, fd)

	g.AddBlocks(gohelper.ModelArray(evt), gohelper.BlockArrayClone(evt), gohelper.ModelArrayToMaps(evt), gohelper.ModelArrayToOrderedMaps(evt))
	if pks := evt.Columns.PKs(); len(pks) > 0 {
		ag, err := gohelper.ModelArrayGet(g, evt, pks, args, goVersion)
		if err != nil {
			return nil, err
		}
		g.AddBlocks(ag)
	}

	return g.Render(linebreak)
}

func eventConstructor(evt *model.Event, x *metamodel.Args) (*golang.Block, error) {
	ret := golang.NewBlock("New"+evt.Proper(), "func")
	argsString, err := helper.GoArgsWithRef(evt.Columns, evt.PackageName(), x)
	if err != nil {
		return nil, err
	}
	ret.WF("func New%s(%s) *%s {", evt.Proper(), argsString, evt.Proper())
	ret.WF("\treturn &%s{%s}", evt.Proper(), evt.Columns.Refs())
	ret.W("}")
	return ret, nil
}

func eventToData(evt *model.Event, cols model.Columns, suffix string, database string) *golang.Block {
	ret := golang.NewBlock(evt.Proper(), "func")
	ret.WF("func (%s *%s) ToData%s() []any {", evt.FirstLetter(), evt.Proper(), suffix)
	calls := lo.Map(cols, func(c *model.Column, _ int) string {
		tk := c.Type.Key()
		complicated := tk == types.KeyAny || tk == types.KeyList || tk == types.KeyMap || tk == types.KeyOrderedMap || tk == types.KeyReference
		if complicated && helper.SimpleJSON(database) {
			return fmt.Sprintf("util.ToJSON(%s.%s),", evt.FirstLetter(), c.Proper())
		} else {
			return fmt.Sprintf("%s.%s,", evt.FirstLetter(), c.Proper())
		}
	})
	lines := util.JoinLines(calls, " ", 120)
	if len(lines) == 1 && len(lines[0]) < 100 {
		ret.WF("\treturn []any{%s}", strings.TrimSuffix(lines[0], ","))
	} else {
		ret.W("\treturn []any{")
		lo.ForEach(lines, func(l string, _ int) {
			ret.WF("\t\t%s", l)
		})
		ret.W("\t}")
	}
	ret.W("}")
	return ret
}

func eventDiff(g *golang.File, evt *model.Event, args *metamodel.Args, linebreak string) (*golang.Block, error) {
	g.AddImport(helper.ImpAppUtil)
	g.AddImport(evt.Imports.Supporting("diff")...)
	return gohelper.DiffBlock(g, evt.Columns, evt, args.Enums)
}
