package goenum

import (
	"fmt"
	"strings"

	"github.com/samber/lo"

	"projectforge.dev/projectforge/app/lib/types"
	"projectforge.dev/projectforge/app/project/export/enum"
	"projectforge.dev/projectforge/app/project/export/golang"
	"projectforge.dev/projectforge/app/util"
)

func enumValues(e *enum.Enum) *golang.Block {
	b := golang.NewBlock(e.Proper(), "vars")
	b.W("var (")

	maxCount := util.StringArrayMaxLength(e.ValuesCamel())
	names := make([]string, 0, len(e.Values))
	pl := len(e.Proper())
	maxColLength := maxCount + pl
	lo.ForEach(e.Values, func(v *enum.Value, _ int) {
		n := e.Proper() + v.Proper()
		names = append(names, n)
		msg := fmt.Sprintf("\t%s = %s{Key: %q", util.StringPad(n, maxColLength), e.Proper(), v.Key)
		if v.Name != "" {
			msg += fmt.Sprintf(", Name: %q", v.Name)
		}
		if v.Description != "" {
			msg += fmt.Sprintf(", Description: %q", v.Description)
		}
		if v.Icon != "" {
			msg += fmt.Sprintf(", Icon: %q", v.Icon)
		}
		ef := e.ExtraFields()
		for _, extraKey := range ef.Order {
			var t string
			if v.Extra != nil {
				t = fmt.Sprint(v.Extra.GetSimple(extraKey))
			}
			switch ef.GetSimple(extraKey) {
			case types.KeyString, "":
				if t == "" {
					continue
				}
				t = "\"" + t + "\""
			case types.KeyTimestamp:
				if t == "<nil>" || t == "" {
					continue
				}
				t = "util.TimeFromStringSimple(\"" + t + "\")"
			case types.KeyBool:
				if t == "<nil>" {
					t = "false"
				}
				if t == "false" || t == "" {
					continue
				}
			}
			msg += fmt.Sprintf(", %s: %s", util.StringToCamel(extraKey), t)
		}

		b.W(msg + "}")
	})

	b.WB()
	b.W("\tAll%s = %s{%s}", e.ProperPlural(), e.ProperPlural(), strings.Join(names, ", "))
	b.W(")")
	return b
}
