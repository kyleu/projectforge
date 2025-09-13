package gohelper

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

func BlockToData(m StringProvider, cols model.Columns, suffix string, database string) *golang.Block {
	ret := golang.NewBlock(m.Proper(), "func")
	ret.WF("func (%s *%s) ToData%s() []any {", m.FirstLetter(), m.Proper(), suffix)
	calls := lo.Map(cols, func(c *model.Column, _ int) string {
		tk := c.Type.Key()
		complicated := tk == types.KeyAny || tk == types.KeyList || tk == types.KeyMap || tk == types.KeyOrderedMap || tk == types.KeyReference
		if complicated && helper.SimpleJSON(database) {
			return fmt.Sprintf("util.ToJSON(%s.%s),", m.FirstLetter(), c.Proper())
		} else {
			return fmt.Sprintf("%s.%s,", m.FirstLetter(), c.Proper())
		}
	})
	lines := util.JoinLines(calls, " ", 120)
	if len(lines) == 1 && len(lines[0]) < 100 {
		ret.WF("\treturn []any{%s}", strings.TrimSuffix(lines[0], ","))
	} else {
		ret.W("\treturn []any{")
		lo.ForEach(lines, func(l string, _ int) {
			ret.WF("\t\t%s", l)
		})
		ret.W("\t}")
	}
	ret.W("}")
	return ret
}
