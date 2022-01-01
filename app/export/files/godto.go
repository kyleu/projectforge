package files

import (
	"fmt"
	"strings"

	"github.com/kyleu/projectforge/app/export/golang"
	"github.com/kyleu/projectforge/app/export/model"
	"github.com/kyleu/projectforge/app/file"
	"github.com/kyleu/projectforge/app/util"
)

func DTO(m *model.Model, args *model.Args) *file.File {
	g := golang.NewFile(m.Package, []string{"app", m.Package}, "dto")
	for _, imp := range importsForTypes("dto", m.Columns.Types()...) {
		g.AddImport(imp.Type, imp.Value)
	}
	g.AddBlocks(modelDTO(m), modelDTOToModel(m), modelDTOArray(), modelDTOArrayTransformer(m))
	return g.Render()
}

func modelDTO(m *model.Model) *golang.Block {
	ret := golang.NewBlock(m.Proper(), "struct")
	ret.W("type dto struct {")
	maxColLength := util.StringArrayMaxLength(m.Columns.CamelNames())
	maxTypeLength := m.Columns.MaxGoDTOKeyLength()
	for _, c := range m.Columns {
		ret.W("\t%s %s `db:%q`", util.StringPad(c.Proper(), maxColLength), util.StringPad(c.ToGoDTOType(), maxTypeLength), c.Name)
	}
	ret.W("}")
	return ret
}

func modelDTOToModel(m *model.Model) *golang.Block {
	ret := golang.NewBlock(m.Proper(), "func")
	ret.W("func (d *dto) To%s() *%s {", m.Proper(), m.Proper())
	refs := make([]string, 0, len(m.Columns))
	for _, c := range m.Columns {
		switch c.Type.Key {
		case model.TypeMap.Key:
			ret.W("\t%s := util.ValueMap{}", c.Camel())
			ret.W("\t_ = util.FromJSON(d.%s, &%s)", c.Proper(), c.Camel())
			refs = append(refs, fmt.Sprintf("%s: %s", c.Proper(), c.Camel()))
		default:
			refs = append(refs, fmt.Sprintf("%s: d.%s", c.Proper(), c.Proper()))
		}
	}
	ret.W("\treturn &%s{%s}", m.Proper(), strings.Join(refs, ", "))
	ret.W("}")
	return ret
}

func modelDTOArray() *golang.Block {
	ret := golang.NewBlock("DTOArray", "type")
	ret.W("type dtos []*dto")
	return ret
}

func modelDTOArrayTransformer(m *model.Model) *golang.Block {
	ret := golang.NewBlock(fmt.Sprintf("DTOTo%s", m.ProperPlural()), "type")
	ret.W("func (x dtos) To%s() %s {", m.ProperPlural(), m.ProperPlural())
	ret.W("\tret := make(%s, 0, len(x))", m.ProperPlural())
	ret.W("\tfor _, d := range x {")
	ret.W("\t\tret = append(ret, d.To%s())", m.Proper())
	ret.W("\t}")
	ret.W("\treturn ret")
	ret.W("}")
	return ret
}
