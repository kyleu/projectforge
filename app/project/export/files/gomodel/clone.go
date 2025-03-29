package gomodel

import (
	"fmt"
	"strings"

	"github.com/samber/lo"

	"projectforge.dev/projectforge/app/lib/metamodel/model"
	"projectforge.dev/projectforge/app/lib/types"
	"projectforge.dev/projectforge/app/project/export/files/helper"
	"projectforge.dev/projectforge/app/project/export/golang"
	"projectforge.dev/projectforge/app/util"
)

func modelClone(g *golang.File, m *model.Model) *golang.Block {
	ret := golang.NewBlock("Clone", "func")
	ret.WF("func (%s *%s) Clone() *%s {", m.FirstLetter(), m.Proper(), m.Proper())
	calls := lo.Map(m.Columns.NotDerived(), func(col *model.Column, _ int) string {
		switch col.Type.Key() {
		case types.KeyMap, types.KeyOrderedMap, types.KeyValueMap, types.KeyReference:
			return fmt.Sprintf("%s.%s.Clone(),", m.FirstLetter(), col.Proper())
		case types.KeyList:
			g.AddImport(helper.ImpAppUtil)
			return fmt.Sprintf("util.ArrayCopy(%s.%s),", m.FirstLetter(), col.Proper())
		default:
			return fmt.Sprintf("%s.%s,", m.FirstLetter(), col.Proper())
		}
	})
	lines := util.JoinLines(calls, " ", 120)
	if len(lines) == 1 && len(lines[0]) < 100 {
		ret.WF("\treturn &%s{%s}", m.Proper(), strings.TrimSuffix(lines[0], ","))
	} else {
		ret.WF("\treturn &%s{", m.Proper())
		lo.ForEach(lines, func(l string, _ int) {
			ret.WF("\t\t%s", l)
		})
		ret.W("\t}")
	}
	ret.W("}")
	return ret
}
