package gomodel

import (
	"github.com/kyleu/projectforge/app/export/golang"
	"github.com/kyleu/projectforge/app/export/model"
)

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
