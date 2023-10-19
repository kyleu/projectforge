package gomodel

import (
	"fmt"
	"strings"

	"github.com/samber/lo"

	"projectforge.dev/projectforge/app/file"
	"projectforge.dev/projectforge/app/lib/types"
	"projectforge.dev/projectforge/app/project/export/enum"
	"projectforge.dev/projectforge/app/project/export/files/helper"
	"projectforge.dev/projectforge/app/project/export/golang"
	"projectforge.dev/projectforge/app/project/export/model"
	"projectforge.dev/projectforge/app/util"
)

func Model(m *model.Model, args *model.Args, addHeader bool, linebreak string) (*file.File, error) {
	g := golang.NewFile(m.Package, []string{"app", m.PackageWithGroup("")}, strings.ToLower(m.Camel()))
	lo.ForEach(helper.ImportsForTypes("go", "", m.Columns.Types()...), func(imp *golang.Import, _ int) {
		g.AddImport(imp)
	})
	lo.ForEach(helper.ImportsForTypes("string", "", m.PKs().Types()...), func(imp *golang.Import, _ int) {
		g.AddImport(imp)
	})
	g.AddImport(helper.ImpAppUtil)
	err := helper.SpecialImports(g, m.Columns, m.PackageWithGroup(""), args.Enums)
	if err != nil {
		return nil, err
	}
	if len(m.PKs()) > 1 {
		pk, e := modelPK(m, args.Enums)
		if e != nil {
			return nil, e
		}
		g.AddBlocks(pk)
	}
	str, err := modelStruct(m, args.Enums)
	if err != nil {
		return nil, err
	}
	c, err := modelConstructor(m, args.Enums)
	if err != nil {
		return nil, err
	}
	rnd, err := modelRandom(m, args.Enums)
	if err != nil {
		return nil, err
	}
	g.AddBlocks(str, c, rnd)
	if b, e := modelFromMap(g, m, args.Enums, args.Database); e == nil {
		g.AddBlocks(b)
	} else {
		return nil, err
	}
	mdiff, err := modelDiff(g, m, args.Enums)
	if err != nil {
		return nil, err
	}

	g.AddBlocks(modelClone(m), modelString(g, m), modelTitle(m))
	if len(m.PKs()) > 1 {
		if pk, e := modelToPK(m, args.Enums); e == nil {
			g.AddBlocks(pk)
		} else {
			return nil, err
		}
	}
	g.AddBlocks(modelWebPath(g, m), mdiff, modelToData(m, m.Columns, "", args.Database))
	if m.IsRevision() {
		hc := m.HistoryColumns(false)
		g.AddBlocks(modelToData(m, hc.Const, "Core", args.Database), modelToData(m, hc.Var, hc.Col.Proper(), args.Database))
	}
	return g.Render(addHeader, linebreak)
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

func modelPK(m *model.Model, enums enum.Enums) (*golang.Block, error) {
	ret := golang.NewBlock("PK", "struct")
	ret.W("type PK struct {")
	pks := m.PKs()
	maxColLength := pks.MaxCamelLength()
	maxTypeLength := pks.MaxGoTypeLength(m.Package, enums)
	for _, c := range pks {
		gt, err := c.ToGoType(m.Package, enums)
		if err != nil {
			return nil, err
		}
		goType := util.StringPad(gt, maxTypeLength)
		ret.W("\t%s %s `json:%q`", util.StringPad(c.Proper(), maxColLength), goType, c.Camel()+modelJSONSuffix(c))
	}
	ret.W("}")
	return ret, nil
}

func modelStruct(m *model.Model, enums enum.Enums) (*golang.Block, error) {
	ret := golang.NewBlock(m.Proper(), "struct")
	ret.W("type %s struct {", m.Proper())
	maxColLength := m.Columns.MaxCamelLength()
	maxTypeLength := m.Columns.MaxGoTypeLength(m.Package, enums)
	for _, c := range m.Columns {
		gt, err := c.ToGoType(m.Package, enums)
		if err != nil {
			return nil, err
		}
		goType := util.StringPad(gt, maxTypeLength)
		tag := fmt.Sprintf("json:\"%s\"", c.Camel()+modelJSONSuffix(c))
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

func modelJSONSuffix(c *model.Column) string {
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
		switch {
		case c.Type.Key() == types.KeyList && database == util.DatabaseSQLite:
			return fmt.Sprintf("util.ToJSON(%s.%s),", m.FirstLetter(), c.Proper())
		default:
			return fmt.Sprintf("%s.%s,", m.FirstLetter(), c.Proper())
		}
	})
	lines := JoinLines(calls, " ", 120)
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
