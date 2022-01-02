package files

import (
	"fmt"
	"strings"

	"github.com/kyleu/projectforge/app/export/golang"
	"github.com/kyleu/projectforge/app/export/model"
	"github.com/kyleu/projectforge/app/file"
	"github.com/kyleu/projectforge/app/util"
)

func Model(m *model.Model, args *model.Args) *file.File {
	g := golang.NewFile(m.Package, []string{"app", m.Package}, m.Camel())
	for _, imp := range importsForTypes("go", m.Columns.Types()...) {
		g.AddImport(imp.Type, imp.Value)
	}
	for _, imp := range importsForTypes("string", m.Columns.PKs().Types()...) {
		g.AddImport(imp.Type, imp.Value)
	}
	if len(m.Columns.PKs()) > 1 {
		g.AddImport(impFmt.Type, impFmt.Value)
	}
	g.AddImport(golang.ImportTypeInternal, "strings")
	g.AddImport(golang.ImportTypeApp, "{{{ .Package }}}/app/util")
	g.AddBlocks(modelStruct(m), modelConstructor(m), modelFromForm(m))
	g.AddBlocks(modelString(m), modelToData(m), modelWebPath(m))
	g.AddBlocks(modelArray(m), modelCols(m))
	return g.Render()
}

func modelCols(m *model.Model) *golang.Block {
	ret := golang.NewBlock("Columns", "procedural")
	ret.W("var (")
	ret.W("\tTable         = %q", m.Name)
	ret.W("\tColumns       = []string{%s}", strings.Join(util.StringArrayQuoted(m.Columns.Names()), ", "))
	ret.W("\tColumnsString = strings.Join(Columns, \", \")")
	ret.W(")")
	return ret
}

func modelStruct(m *model.Model) *golang.Block {
	ret := golang.NewBlock(m.Proper(), "struct")
	ret.W("type %s struct {", m.Proper())
	maxColLength := util.StringArrayMaxLength(m.Columns.CamelNames())
	maxTypeLength := m.Columns.MaxGoKeyLength()
	for _, c := range m.Columns {
		ret.W("\t%s %s `json:%q`", util.StringPad(c.Proper(), maxColLength), util.StringPad(c.ToGoType(), maxTypeLength), c.Camel())
	}
	ret.W("}")
	return ret
}

func modelConstructor(m *model.Model) *golang.Block {
	ret := golang.NewBlock("New"+m.Proper(), "func")
	ret.W("func New(%s) *%s {", m.Columns.PKs().Args(), m.Proper())
	ret.W("\treturn &%s{%s}", m.Proper(), m.Columns.PKs().Refs())
	ret.W("}")
	return ret
}

func modelFromForm(m *model.Model) *golang.Block {
	ret := golang.NewBlock(m.Package+"FromForm", "func")
	ret.W("func FromMap(m util.ValueMap, setPK bool) (*%s, error) {", m.Proper())
	ret.W("\tret := &%s{}", m.Proper())
	ret.W("\tvar err error")
	ret.W("\tif setPK {")
	for _, col := range m.Columns.PKs() {
		ret.W("\t\tret.%s, err = m.Parse%s(%q)", col.Proper(), col.ToGoMapParse(), col.Camel())
		ret.W("\t\tif err != nil {")
		ret.W("\t\t\treturn nil, err")
		ret.W("\t\t}")
	}
	ret.W("\t}")
	for _, col := range m.Columns {
		if !col.PK {
			ret.W("\tret.%s, err = m.Parse%s(%q)", col.Proper(), col.ToGoMapParse(), col.Camel())
			ret.W("\tif err != nil {")
			ret.W("\t\treturn nil, err")
			ret.W("\t}")
		}
	}
	ret.W("\treturn ret, nil")
	ret.W("}")

	return ret
}

func modelString(m *model.Model) *golang.Block {
	ret := golang.NewBlock("String", "func")
	ret.W("func (%s *%s) String() string {", m.FirstLetter(), m.Proper())
	if pks := m.Columns.PKs(); len(pks) == 1 {
		switch pks[0].Type.Key {
		case model.TypeString.Key:
			ret.W("\treturn %s.%s", m.FirstLetter(), pks[0].Proper())
		case model.TypeUUID.Key:
			ret.W("\treturn %s.%s.String()", m.FirstLetter(), pks[0].Proper())
		default:
			ret.W("\treturn fmt.Sprint(%s.%s)", m.FirstLetter(), pks[0].Proper())
		}
	} else {
		s := "\treturn fmt.Sprintf(\""
		for idx := range m.Columns.PKs() {
			if idx > 0 {
				s += "::"
			}
			s += "%%s"
		}
		s += "\""
		for _, c := range m.Columns.PKs() {
			s += ", " + c.ToGoString(m.FirstLetter()+".")
		}
		ret.W(s + ")")
	}
	ret.W("}")
	return ret
}

func modelToData(m *model.Model) *golang.Block {
	ret := golang.NewBlock(m.Proper(), "func")
	ret.W("func (%s *%s) ToData() []interface{} {", m.FirstLetter(), m.Proper())
	refs := make([]string, 0, len(m.Columns))
	for _, c := range m.Columns {
		refs = append(refs, fmt.Sprintf("%s.%s", m.FirstLetter(), c.Proper()))
	}
	ret.W("\treturn []interface{}{%s}", strings.Join(refs, ", "))
	ret.W("}")
	return ret
}

func modelWebPath(m *model.Model) *golang.Block {
	ret := golang.NewBlock("WebPath", "type")
	ret.W("func (%s *%s) WebPath() string {", m.FirstLetter(), m.Proper())
	p := "\"/" + m.Package + "\""
	for _, pk := range m.Columns.PKs() {
		p += " + \"/\" + "
		p += pk.ToGoString(m.FirstLetter() + ".")
	}
	ret.W("\treturn " + p)
	ret.W("}")
	return ret
}

func modelArray(m *model.Model) *golang.Block {
	ret := golang.NewBlock(m.Proper()+"Array", "type")
	ret.W("type %s []*%s", m.ProperPlural(), m.Proper())
	return ret
}
