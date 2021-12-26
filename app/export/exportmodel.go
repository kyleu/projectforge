package export

import (
	"fmt"
	"strings"

	"github.com/kyleu/projectforge/app/file"
	"github.com/kyleu/projectforge/app/util"
)

func exportModelFile(m *Model, args *Args) *file.File {
	g := NewGoFile(m.Package, []string{"app", m.Package}, m.camel())
	for _, imp := range m.Columns.Types().Imports() {
		g.AddImport(imp.Type, imp.Value)
	}
	g.AddImport(ImportTypeInternal, "fmt")
	g.AddImport(ImportTypeInternal, "strings")
	g.AddBlocks(modelStruct(m), modelConstructor(m), modelString(m), modelToData(m), modelArray(m), modelCols(m))
	g.AddBlocks(modelDTO(m), modelDTOToModel(m), modelDTOArray(), modelDTOArrayTransformer(m))
	return g.Render()
}

func modelCols(m *Model) *Block {
	ret := NewBlock("Columns", "procedural")
	ret.W("var (")
	ret.WF("\tTable         = %q", m.Name)
	ret.WF("\tColumns       = []string{%s}", strings.Join(util.StringArrayQuoted(m.Columns.Names()), ", "))
	ret.W("\tColumnsString = strings.Join(Columns, \", \")")
	ret.W(")")
	return ret
}

func modelStruct(m *Model) *Block {
	ret := NewBlock(m.proper(), "struct")
	ret.WF("type %s struct {", m.proper())
	maxColLength := util.StringArrayMaxLength(m.Columns.camelNames())
	maxTypeLength := m.Columns.Types().MaxGoKeyLength()
	for _, c := range m.Columns {
		ret.WF("\t%s %s `json:%q`", util.StringPad(c.proper(), maxColLength), util.StringPad(c.Type.Go, maxTypeLength), c.camel())
	}
	ret.W("}")
	return ret
}

func modelConstructor(m *Model) *Block {
	ret := NewBlock("New"+m.proper(), "func")
	ret.WF("func New(%s) *%s {", m.Columns.PKs().Args(), m.proper())
	ret.WF("\treturn &%s{%s}", m.proper(), m.Columns.PKs().Refs())
	ret.W("}")
	return ret
}

func modelString(m *Model) *Block {
	ret := NewBlock("String", "func")
	ret.WF("func (%s *%s) String() string {", m.firstLetter(), m.proper())
	if pks := m.Columns.PKs(); len(pks) == 1 {
		ret.WF("\treturn fmt.Sprint(%s.%s)", m.firstLetter(), pks[0].proper())
	} else {
		s := "\treturn fmt.Sprintf(\""
		for idx := range m.Columns.PKs() {
			if idx > 0 {
				s += "::"
			}
			s += "%%%%s"
		}
		s += "\""
		for _, c := range m.Columns.PKs() {
			s += fmt.Sprintf(", %s.%s", m.firstLetter(), c.proper())
		}
		ret.WF(s + ")")
	}
	ret.W("}")
	return ret
}

func modelToData(m *Model) *Block {
	ret := NewBlock(m.proper(), "func")
	ret.WF("func (%s *%s) ToData() []interface{} {", m.firstLetter(), m.proper())
	refs := make([]string, 0, len(m.Columns))
	for _, c := range m.Columns {
		refs = append(refs, fmt.Sprintf("%s.%s", m.firstLetter(), c.proper()))
	}
	ret.WF("\treturn []interface{}{%s}", strings.Join(refs, ", "))
	ret.W("}")
	return ret
}

func modelArray(m *Model) *Block {
	ret := NewBlock(m.proper()+"Array", "type")
	ret.WF("type %s []*%s", m.properPlural(), m.proper())
	return ret
}

func modelDTO(m *Model) *Block {
	ret := NewBlock(m.proper(), "struct")
	ret.W("type dto struct {")
	maxColLength := util.StringArrayMaxLength(m.Columns.camelNames())
	maxTypeLength := m.Columns.Types().MaxGoKeyLength()
	for _, c := range m.Columns {
		ret.WF("\t%s %s `db:%q`", util.StringPad(c.proper(), maxColLength), util.StringPad(c.Type.Go, maxTypeLength), c.Name)
	}
	ret.W("}")
	return ret
}

func modelDTOToModel(m *Model) *Block {
	ret := NewBlock(m.proper(), "func")
	ret.WF("func (d *dto) To%s() *%s {", m.proper(), m.proper())
	refs := make([]string, 0, len(m.Columns))
	for _, c := range m.Columns {
		refs = append(refs, fmt.Sprintf("%s: d.%s", c.proper(), c.proper()))
	}
	ret.WF("\treturn &%s{%s}", m.proper(), strings.Join(refs, ", "))
	ret.W("}")
	return ret
}

func modelDTOArray() *Block {
	ret := NewBlock("DTOArray", "type")
	ret.W("type dtos []*dto")
	return ret
}

func modelDTOArrayTransformer(m *Model) *Block {
	ret := NewBlock(fmt.Sprintf("DTOTo%s", m.properPlural()), "type")
	ret.WF("func (x dtos) To%s() %s {", m.properPlural(), m.properPlural())
	ret.WF("\tret := make(%s, 0, len(x))", m.properPlural())
	ret.W("\tfor _, d := range x {")
	ret.WF("\t\tret = append(ret, d.To%s())", m.proper())
	ret.W("\t}")
	ret.W("\treturn ret")
	ret.W("}")
	return ret
}
