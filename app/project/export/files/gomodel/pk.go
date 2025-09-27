package gomodel

import (
	"projectforge.dev/projectforge/app/lib/metamodel/enum"
	"projectforge.dev/projectforge/app/lib/metamodel/model"
	"projectforge.dev/projectforge/app/lib/types"
	"projectforge.dev/projectforge/app/project/export/files/gohelper"
	"projectforge.dev/projectforge/app/project/export/golang"
	"projectforge.dev/projectforge/app/util"
)

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
		ret.WF("\t%s %s %s", util.StringPad(c.Proper(), maxColLength), goType, gohelper.ColumnTag(c))
	}
	ret.W("}")
	return ret, nil
}

func modelPKString(m *model.Model) (*golang.Block, error) {
	ret := golang.NewBlock("PKString", "func")
	ret.W("func (p *PK) String() string {")
	pks := m.PKs()
	format, args := make([]string, 0, len(pks)), make([]string, 0, len(pks))
	for _, c := range pks {
		switch c.Type.Key() {
		case types.KeyString, "":
			format = append(format, "%%s")
		default:
			format = append(format, "%%v")
		}
		args = append(args, "p."+c.Proper())
	}
	ret.WF("\treturn fmt.Sprintf(%q, %s)", util.StringJoin(format, " • "), util.StringJoin(args, ", "))
	ret.W("}")
	return ret, nil
}
