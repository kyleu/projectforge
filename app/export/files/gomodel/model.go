package gomodel

import (
	"fmt"
	"strings"

	"github.com/kyleu/projectforge/app/export/files/helper"
	"github.com/kyleu/projectforge/app/export/golang"
	"github.com/kyleu/projectforge/app/export/model"
	"github.com/kyleu/projectforge/app/file"
	"github.com/kyleu/projectforge/app/util"
)

func Model(m *model.Model, args *model.Args) (*file.File, error) {
	g := golang.NewFile(m.Package, []string{"app", m.Package}, m.Camel())
	for _, imp := range helper.ImportsForTypes("go", m.Columns.Types()...) {
		g.AddImport(imp)
	}
	for _, imp := range helper.ImportsForTypes("string", m.PKs().Types()...) {
		g.AddImport(imp)
	}
	if len(m.PKs()) > 1 {
		g.AddImport(helper.ImpFmt)
	}
	g.AddImport(helper.ImpAppUtil)
	g.AddBlocks(modelStruct(m), modelConstructor(m), modelFromMap(m))
	g.AddBlocks(modelString(m), modelWebPath(m), modelToData(m, m.Columns, ""))
	if m.IsRevision() {
		hc := m.HistoryColumns(false)
		g.AddBlocks(modelToData(m, hc.Const, "Core"), modelToData(m, hc.Var, hc.Col.Proper()))
	}
	g.AddBlocks(modelArray(m))
	return g.Render()
}

func modelStruct(m *model.Model) *golang.Block {
	ret := golang.NewBlock(m.Proper(), "struct")
	ret.W("type %s struct {", m.Proper())
	maxColLength := util.StringArrayMaxLength(m.Columns.CamelNames())
	maxTypeLength := m.Columns.MaxGoKeyLength()
	for _, c := range m.Columns {
		suffix := ""
		if c.Nullable {
			suffix = ",omitempty"
		}
		ret.W("\t%s %s `json:%q`", util.StringPad(c.Proper(), maxColLength), util.StringPad(c.ToGoType(), maxTypeLength), c.Camel()+suffix)
	}
	ret.W("}")
	return ret
}

func modelConstructor(m *model.Model) *golang.Block {
	ret := golang.NewBlock("New"+m.Proper(), "func")
	ret.W("func New(%s) *%s {", m.PKs().Args(), m.Proper())
	ret.W("\treturn &%s{%s}", m.Proper(), m.PKs().Refs())
	ret.W("}")
	return ret
}

func modelFromMap(m *model.Model) *golang.Block {
	ret := golang.NewBlock(m.Package+"FromForm", "func")
	ret.W("func FromMap(m util.ValueMap, setPK bool) (*%s, error) {", m.Proper())
	ret.W("\tret := &%s{}", m.Proper())
	ret.W("\tvar err error")
	ret.W("\tif setPK {")
	cols := m.Columns.WithoutTag("created").WithoutTag("updated").WithoutTag(model.RevisionType)
	for _, col := range cols.PKs() {
		ret.W("\t\tret.%s, err = m.Parse%s(%q)", col.Proper(), col.ToGoMapParse(), col.Camel())
		ret.W("\t\tif err != nil {")
		ret.W("\t\t\treturn nil, err")
		ret.W("\t\t}")
	}
	ret.W("\t}")
	for _, col := range cols {
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
	if pks := m.PKs(); len(pks) == 1 {
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
		for idx := range m.PKs() {
			if idx > 0 {
				s += "::"
			}
			s += "%%s"
		}
		s += "\""
		for _, c := range m.PKs() {
			s += ", " + c.ToGoString(m.FirstLetter()+".")
		}
		ret.W(s + ")")
	}
	ret.W("}")
	return ret
}

func modelToData(m *model.Model, cols model.Columns, suffix string) *golang.Block {
	ret := golang.NewBlock(m.Proper(), "func")
	ret.W("func (%s *%s) ToData%s() []interface{} {", m.FirstLetter(), m.Proper(), suffix)
	refs := make([]string, 0, len(cols))
	for _, c := range cols {
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
	for _, pk := range m.PKs() {
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
