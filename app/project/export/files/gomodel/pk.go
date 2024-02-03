package gomodel

import (
	"strings"

	"projectforge.dev/projectforge/app/project/export/enum"
	"projectforge.dev/projectforge/app/project/export/golang"
	"projectforge.dev/projectforge/app/project/export/model"
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
		ret.W("\t%s %s `json:%q`", util.StringPad(c.Proper(), maxColLength), goType, c.Camel()+modelJSONSuffix(c))
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
		case "string", "":
			format = append(format, "%%s")
		default:
			format = append(format, "%%v")
		}
		args = append(args, "p."+c.Proper())
	}
	ret.W("\treturn fmt.Sprintf(%q, %s)", strings.Join(format, "::"), strings.Join(args, ", "))
	ret.W("}")
	return ret, nil
}
