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
	if !m.SkipController() {
		g.AddImport(helper.ImpPath)
	}
	g.AddImport(helper.ImpAppUtil, helper.ImpAppSvc)
	imps, err := helper.SpecialImports(m.Columns, m.PackageWithGroup(""), args.Models, args.Enums, args.ExtraTypes)
	if err != nil {
		return nil, err
	}
	g.AddImport(imps...)
	imps, err = helper.EnumImports(m.Columns.Types(), m.PackageWithGroup(""), args.Models, args.Enums)
	if err != nil {
		return nil, err
	}
	g.AddImport(imps...)
	g.AddImport(m.Imports.Supporting("model")...)

	if !m.SkipController() {
		g.AddBlocks(defaultRoute(m), routeMethod(m))
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
	str, err := modelStruct(m, args.Models, args.Enums, args.ExtraTypes)
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
	if !m.SkipController() {
		g.AddBlocks(modelWebPath(g, m))
	}
	g.AddBlocks(modelToData(m, m.Columns.NotDerived(), "", args.Database), fd)
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

func modelStruct(m *model.Model, models model.Models, enums enum.Enums, extraTypes model.Models) (*golang.Block, error) {
	ret := golang.NewBlock(m.Proper(), "struct")
	ret.WF("type %s struct {", m.Proper())
	cols := m.Columns.NotDerived()
	maxColLength := cols.MaxCamelLength()
	maxTypeLength := cols.MaxGoTypeLength(m.Package, enums)

	gts := lo.Map(cols, func(c *model.Column, _ int) string {
		gt, err := c.ToGoType(m.Package, enums)
		if err != nil {
			return err.Error()
		}
		if ref, mdl, _ := helper.LoadRef(c, models, extraTypes); ref != nil && !strings.Contains(gt, ".") {
			if mdl != nil && mdl.Package != m.Package {
				gt = mdl.Package + "." + gt
				if x := len(gt); maxTypeLength < x {
					maxTypeLength = x
				}
			}
		}
		if gt == "*"+types.KeyAny {
			gt = types.KeyAny
		}
		return gt
	})

	for idx, c := range cols {
		goType := util.StringPad(gts[idx], maxTypeLength)
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
		ret.WF("\t%s %s `%s`", util.StringPad(c.Proper(), maxColLength), goType, tag)
	}
	ret.W("}")
	return ret, nil
}

func modelJSONSuffix(col *model.Column) string {
	if col.HasTag("force-json") {
		return ""
	}
	return ",omitempty"
}

func modelConstructor(m *model.Model, enums enum.Enums) (*golang.Block, error) {
	ret := golang.NewBlock("New"+m.Proper(), "func")
	args, err := m.PKs().Args(m.Package, enums)
	if err != nil {
		return nil, err
	}
	ret.WF("func New%s(%s) *%s {", m.Proper(), args, m.Proper())
	ret.WF("\treturn &%s{%s}", m.Proper(), m.PKs().Refs())
	ret.W("}")
	return ret, nil
}

func modelToData(m *model.Model, cols model.Columns, suffix string, database string) *golang.Block {
	ret := golang.NewBlock(m.Proper(), "func")
	ret.WF("func (%s *%s) ToData%s() []any {", m.FirstLetter(), m.Proper(), suffix)
	calls := lo.Map(cols, func(c *model.Column, _ int) string {
		tk := c.Type.Key()
		complicated := tk == types.KeyAny || tk == types.KeyList || tk == types.KeyMap || tk == types.KeyOrderedMap || tk == types.KeyReference
		if complicated && helper.SimpleJSON(database) {
			return fmt.Sprintf("util.ToJSON(%s.%s),", m.FirstLetter(), c.Proper())
		} else {
			return fmt.Sprintf("%s.%s,", m.FirstLetter(), c.Proper())
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
