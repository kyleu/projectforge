package gomodel

import (
	"projectforge.dev/projectforge/app/lib/metamodel/model"
	"projectforge.dev/projectforge/app/project/export/golang"
)

func modelToCSV(m *model.Model) *golang.Block {
	ret := golang.NewBlock("ToCSV", "func")
	ret.WF("func (%s *%s) ToCSV() ([]string, [][]string) {", m.FirstLetter(), m.Proper())
	ret.WF("\treturn %sFieldDescs.Keys(), [][]string{%s.Strings()}", m.Proper(), m.FirstLetter())
	ret.W("}")
	return ret
}

func modelArrayToCSV(m *model.Model) *golang.Block {
	ret := golang.NewBlock("ToCSV", "func")
	ret.WF("func (%s %s) ToCSV() ([]string, [][]string) {", m.FirstLetter(), m.ProperPlural())
	ret.WF("\treturn %sFieldDescs.Keys(), lo.Map(%s, func(x *%s, _ int) []string {", m.Proper(), m.FirstLetter(), m.Proper())
	ret.W("\t\treturn x.Strings()")
	ret.W("\t})")
	ret.W("}")
	return ret
}
