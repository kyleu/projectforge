package gohelper

import (
	"fmt"
	"strings"

	"github.com/samber/lo"

	"projectforge.dev/projectforge/app/lib/metamodel"
	"projectforge.dev/projectforge/app/lib/metamodel/model"
	"projectforge.dev/projectforge/app/lib/types"
	"projectforge.dev/projectforge/app/project/export/files/helper"
	"projectforge.dev/projectforge/app/project/export/golang"
	"projectforge.dev/projectforge/app/util"
)

func BlockClone(g *golang.File, cols model.Columns, str metamodel.StringProvider) *golang.Block {
	ret := golang.NewBlock("Clone", "func")
	ret.WF("func (%s *%s) Clone() *%s {", str.FirstLetter(), str.Proper(), str.Proper())
	ret.WF("\tif %s == nil {", str.FirstLetter())
	ret.WF("\t\treturn nil")
	ret.WF("\t}")
	calls := lo.Map(cols, func(col *model.Column, _ int) string {
		switch col.Type.Key() {
		case types.KeyMap, types.KeyOrderedMap, types.KeyValueMap, types.KeyNumericMap, types.KeyReference:
			return fmt.Sprintf("%s: %s.%s.Clone(),", col.Proper(), str.FirstLetter(), col.Proper())
		case types.KeyList:
			g.AddImport(helper.ImpAppUtil)
			return fmt.Sprintf("%s: util.ArrayCopy(%s.%s),", col.Proper(), str.FirstLetter(), col.Proper())
		default:
			return fmt.Sprintf("%s: %s.%s,", col.Proper(), str.FirstLetter(), col.Proper())
		}
	})
	lines := util.JoinLines(calls, " ", 120)
	if len(lines) == 1 && len(lines[0]) < (110-len(str.Proper())) {
		ret.WF("\treturn &%s{%s}", str.Proper(), strings.TrimSuffix(lines[0], ","))
	} else {
		ret.WF("\treturn &%s{", str.Proper())
		lo.ForEach(lines, func(l string, _ int) {
			ret.WF("\t\t%s", l)
		})
		ret.W("\t}")
	}
	ret.W("}")
	return ret
}

func BlockArrayClone(str metamodel.StringProvider) *golang.Block {
	ret := golang.NewBlock(str.Proper()+"ArrayClone", "func")
	ret.WF("func (%s %s) Clone() %s {", str.FirstLetter(), str.ProperPlural(), str.ProperPlural())
	ret.WF("\treturn lo.Map(%s, func(xx *%s, _ int) *%s {", str.FirstLetter(), str.Proper(), str.Proper())
	ret.W("\t\treturn xx.Clone()")
	ret.W("\t})")
	ret.W("}")
	return ret
}
