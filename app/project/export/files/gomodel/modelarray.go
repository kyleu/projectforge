package gomodel

import (
	"fmt"
	"projectforge.dev/projectforge/app/project/export/enum"
	"strings"

	"projectforge.dev/projectforge/app/project/export/golang"
	"projectforge.dev/projectforge/app/project/export/model"
)

func modelArray(m *model.Model) *golang.Block {
	ret := golang.NewBlock(m.Proper()+"Array", "type")
	ret.W("type %s []*%s", m.ProperPlural(), m.Proper())
	return ret
}

func modelArrayGet(m *model.Model, enums enum.Enums) (*golang.Block, error) {
	ret := golang.NewBlock(m.Proper()+"ArrayGet", "func")
	args, err := m.PKs().Args(m.Package, enums)
	if err != nil {
		return nil, err
	}
	ret.W("func (%s %s) Get(%s) *%s {", m.FirstLetter(), m.ProperPlural(), args, m.Proper())
	ret.W("\tfor _, x := range %s {", m.FirstLetter())
	comps := make([]string, 0, len(m.PKs()))
	for _, pk := range m.PKs() {
		comps = append(comps, fmt.Sprintf("x.%s == %s", pk.Proper(), pk.Camel()))
	}
	ret.W("\t\tif %s {", strings.Join(comps, " && "))
	ret.W("\t\t\treturn x")
	ret.W("\t\t}")
	ret.W("\t}")
	ret.W("\treturn nil")
	ret.W("}")
	return ret, nil
}

func modelArrayClone(m *model.Model) *golang.Block {
	ret := golang.NewBlock(m.Proper()+"ArrayClone", "func")
	ret.W("func (%s %s) Clone() %s {", m.FirstLetter(), m.ProperPlural(), m.ProperPlural())
	ret.W("\treturn slices.Clone(%s)", m.FirstLetter())
	ret.W("}")
	return ret
}
