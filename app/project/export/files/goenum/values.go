package goenum

import (
	"fmt"
	"strings"

	"github.com/samber/lo"

	"projectforge.dev/projectforge/app/project/export/enum"
	"projectforge.dev/projectforge/app/project/export/golang"
	"projectforge.dev/projectforge/app/util"
)

func enumValues(e *enum.Enum, g *golang.File) *golang.Block {
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
		for extraKey, extraType := range e.ExtraFields() {
			t := v.Extra.GetStringOpt(extraKey)
			if extraType == "string" || extraType == "" {
				t = "\"" + t + "\""
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
