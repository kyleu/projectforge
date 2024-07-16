package util

import (
	"strings"
	"unicode"

	"github.com/pkg/errors"
	"github.com/samber/lo"
)

var acronymMap = map[string]string{}

func ConfigureAcronym(key, val string) {
	acronymMap[key] = val
}

var acronyms []string

func InitAcronyms(extras ...string) error {
	if len(acronyms) > 0 {
		return errors.New("double initialization of acronyms")
	}
	x := []string{"Api", "Html", "Id", "Ip", "Json", "Md5", "Sha", "Sku", "Sql", "Xml", "Uri", "Url"}
	x = append(x, lo.Map(extras, func(s string, _ int) string {
		return strings.ToUpper(s[:1]) + strings.ToLower(s[1:])
	})...)
	lo.ForEach(x, func(x string, _ int) {
		ConfigureAcronym(strings.ToUpper(x), strings.ToLower(x))
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
