package gohelper

import (
	"projectforge.dev/projectforge/app/lib/metamodel"
	"projectforge.dev/projectforge/app/project/export/golang"
)

func BlockToCSV(str metamodel.StringProvider) *golang.Block {
	ret := golang.NewBlock("ToCSV", "func")
	ret.WF("func (%s *%s) ToCSV() ([]string, [][]string) {", str.FirstLetter(), str.Proper())
	ret.WF("\treturn %sFieldDescs.Keys(), [][]string{%s.Strings()}", str.Proper(), str.FirstLetter())
	ret.W("}")
	return ret
}

func BlockArrayToCSV(str metamodel.StringProvider) *golang.Block {
	ret := golang.NewBlock("ToCSV", "func")
	ret.WF("func (%s %s) ToCSV() ([]string, [][]string) {", str.FirstLetter(), str.ProperPlural())
	ret.WF("\treturn %sFieldDescs.Keys(), lo.Map(%s, func(x *%s, _ int) []string {", str.Proper(), str.FirstLetter(), str.Proper())
	ret.W("\t\treturn x.Strings()")
	ret.W("\t})")
	ret.W("}")
	return ret
}
