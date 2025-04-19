package goenum

import (
	"fmt"
	"strings"

	"github.com/samber/lo"

	"projectforge.dev/projectforge/app/lib/metamodel/enum"
	"projectforge.dev/projectforge/app/lib/types"
	"projectforge.dev/projectforge/app/project/export/files/helper"
	"projectforge.dev/projectforge/app/project/export/golang"
	"projectforge.dev/projectforge/app/util"
)

func enumValues(e *enum.Enum) *golang.Block {
	names := lo.Map(e.Values, func(x *enum.Value, _ int) string {
		return e.Proper() + x.Proper(e.Acronyms...)
	})

	b := golang.NewBlock(e.Proper(), "vars")

	if e.Simple() {
		b.WF("var All%s = %s{%s}", e.ProperPlural(), e.ProperPlural(), strings.Join(names, ", "))
		return b
	}

	b.W("var (")
	maxCount := util.StringArrayMaxLength(e.ValuesCamel())

	pl := len(e.Proper())
	maxColLength := maxCount + pl
	lo.ForEach(e.Values, func(v *enum.Value, _ int) {
		b.W(enumValue(e, v, maxColLength))
	})

	b.WB()
	b.WF("\tAll%s = %s{%s}", e.ProperPlural(), e.ProperPlural(), strings.Join(names, ", "))
	b.W(")")
	return b
}

func enumValue(e *enum.Enum, v *enum.Value, maxColLength int) string {
	msg := fmt.Sprintf("\t%s = %s{Key: %q", util.StringPad(e.Proper()+v.Proper(e.Acronyms...), maxColLength), e.Proper(), v.Key)
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
		case types.KeyTimestamp, types.KeyTimestampZoned:
			if t == helper.TextNil || t == "" {
				continue
			}
			t = "util.TimeFromStringSimple(\"" + t + "\")"
		case types.KeyBool:
			if t == helper.TextNil {
				t = util.BoolFalse
			}
			if t == util.BoolFalse || t == "" {
				continue
			}
		}
		msg += fmt.Sprintf(", %s: %s", util.StringToProper(extraKey), t)
	}
	msg += "}"
	return msg
}
