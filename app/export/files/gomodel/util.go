package gomodel

import (
	"github.com/kyleu/projectforge/app/export/golang"
	"github.com/kyleu/projectforge/app/export/model"
	"github.com/kyleu/projectforge/app/lib/types"
	"github.com/kyleu/projectforge/app/util"
)

func modelClone(m *model.Model) *golang.Block {
	ret := golang.NewBlock("Clone", "func")
	ret.W("func (%s *%s) Clone() *%s {", m.FirstLetter(), m.Proper(), m.Proper())
	ret.W("\treturn &%s{", m.Proper())
	max := m.Columns.MaxCamelLength() + 1
	for _, col := range m.Columns {
		ret.W("\t\t%s %s.%s,", util.StringPad(col.Proper()+":", max), m.FirstLetter(), col.Proper())
	}
	ret.W("\t}")
	ret.W("}")
	return ret
}

func modelString(m *model.Model) *golang.Block {
	ret := golang.NewBlock("String", "func")
	ret.W("func (%s *%s) String() string {", m.FirstLetter(), m.Proper())
	if pks := m.PKs(); len(pks) == 1 {
		switch pks[0].Type.Key() {
		case types.KeyString:
			ret.W("\treturn %s.%s", m.FirstLetter(), pks[0].Proper())
		case types.KeyUUID:
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

func modelWebPath(m *model.Model) *golang.Block {
	ret := golang.NewBlock("WebPath", "type")
	ret.W("func (%s *%s) WebPath() string {", m.FirstLetter(), m.Proper())
	p := "\"/" + m.Route() + "\""
	for _, pk := range m.PKs() {
		p += " + \"/\" + "
		p += pk.ToGoString(m.FirstLetter() + ".")
	}
	ret.W("\treturn " + p)
	ret.W("}")
	return ret
}
