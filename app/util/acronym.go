// Package util - Content managed by Project Forge, see [projectforge.md] for details.
package util

import (
	"strings"
	"unicode"

	"github.com/iancoleman/strcase"
	"github.com/pkg/errors"
	"github.com/samber/lo"
)

func StringToCamel(s string, extraAcronyms ...string) string {
	return acr(strcase.ToCamel(s), extraAcronyms...)
}

func StringToTitle(s string, extraAcronyms ...string) string {
	ret := strings.Builder{}
	runes := []rune(StringToCamel(s, extraAcronyms...))
	lo.ForEach(runes, func(c rune, idx int) {
		if idx > 0 && idx < len(runes)-1 && unicode.IsUpper(c) {
			if !unicode.IsUpper(runes[idx+1]) {
				ret.WriteRune(' ')
			} else if !unicode.IsUpper(runes[idx-1]) {
				ret.WriteRune(' ')
			}
		}
		ret.WriteRune(c)
	})
	return ret.String()
}

func StringToLowerCamel(s string, extraAcronyms ...string) string {
	return acr(strcase.ToLowerCamel(s), extraAcronyms...)
}

func StringToSnake(s string, extraAcronyms ...string) string {
	return acr(strcase.ToSnake(s), extraAcronyms...)
}

var acronyms []string

func InitAcronyms(extras ...string) error {
	if len(acronyms) > 0 {
		return errors.New("double initialization of acronyms")
	}
	x := []string{"Api", "Html", "Id", "Ip", "Json", "Sql", "Xml", "Uri", "Url"}
	x = append(x, lo.Map(extras, func(s string, _ int) string {
		return strings.ToUpper(s[:1]) + strings.ToLower(s[1:])
	})...)
	lo.ForEach(x, func(x string, _ int) {
		strcase.ConfigureAcronym(strings.ToUpper(x), strings.ToLower(x))
	})
	acronyms = x
	return nil
}

func acr(ret string, extraAcronyms ...string) string {
	proc := func(a string) {
		var lastIdx int
		for {
			i := strings.Index(ret[lastIdx:], a)
			if i == -1 {
				if strings.EqualFold(a, ret) && unicode.IsUpper(rune(ret[0])) {
					ret = strings.ToUpper(a)
				}
				break
			}
			i += lastIdx
			lastIdx = i + len(a)
			if lastIdx >= len(ret) {
				ret = ret[:i] + strings.ToUpper(a) + ret[lastIdx:]
			} else {
				s := string(ret[lastIdx])
				if strings.ToUpper(s) == s {
					ret = ret[:i] + strings.ToUpper(a) + ret[lastIdx:]
				} else {
					lastIdx++
				}
			}
		}
	}
	lo.ForEach(acronyms, func(a string, _ int) {
		proc(a)
	})
	lo.ForEach(extraAcronyms, func(a string, _ int) {
		proc(a)
	})
	return ret
}
