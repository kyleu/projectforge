package gomodel

import (
	"fmt"
	"strings"

	"projectforge.dev/projectforge/app/export/golang"
	"projectforge.dev/projectforge/app/export/model"
)

func modelArray(m *model.Model) *golang.Block {
	ret := golang.NewBlock(m.Proper()+"Array", "type")
	ret.W("type %s []*%s", m.ProperPlural(), m.Proper())
	return ret
}

func modelArrayGet(m *model.Model) *golang.Block {
	ret := golang.NewBlock(m.Proper()+"ArrayGet", "func")
	ret.W("func (%s %s) Get(%s) *%s {", m.FirstLetter(), m.ProperPlural(), m.PKs().Args(m.Package), m.Proper())
	ret.W("\tfor _, x := range %s {", m.FirstLetter())
	var comps []string
	for _, pk := range m.PKs() {
		comps = append(comps, fmt.Sprintf("x.%s == %s", pk.Proper(), pk.Camel()))
	}
	ret.W("\t\tif %s {", strings.Join(comps, " && "))
	ret.W("\t\t\treturn x")
	ret.W("\t\t}")
	ret.W("\t}")
	ret.W("\treturn nil")
	ret.W("}")
	return ret
}

func modelArrayClone(m *model.Model) *golang.Block {
	ret := golang.NewBlock(m.Proper()+"ArrayClone", "func")
	ret.W("func (%s %s) Clone() %s {", m.FirstLetter(), m.ProperPlural(), m.ProperPlural())
	ret.W("\treturn slices.Clone(%s)", m.FirstLetter())
	ret.W("}")
	return ret
}
