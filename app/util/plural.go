package util

import (
	"fmt"
	"github.com/samber/lo"
	"strings"
	"unicode"

	"github.com/gertd/go-pluralize"
)

var plrl *pluralize.Client

func plrlSvc() {
	if plrl == nil {
		plrl = pluralize.NewClient()
	}
}

func StringToPlural(s string) string {
	plrlSvc()
	ret := plrl.Plural(s)
	if len(ret) < 3 || len(ret) == len(s) {
		return ret
	}
	runes := []rune(ret)
	l0, l1, l2 := ret[len(ret)-1], runes[len(runes)-2], runes[len(runes)-3]
	if l0 == 'S' && unicode.IsUpper(l1) && unicode.IsUpper(l2) {
		runes[len(runes)-1] = 's'
		ret = string(runes)
	}
	if l1 == 'E' && unicode.IsUpper(l2) {
		runes[len(runes)-2] = 'e'
		ret = string(runes)
	}
	return ret
}

func StringToSingular(s string) string {
	plrlSvc()
	ret := plrl.Singular(s)
	if len(s) != len(ret) {
		orig := s[:len(ret)]
		if lo.EveryBy([]rune(orig), func(x rune) bool {
			return unicode.IsUpper(x)
		}) && strings.ToLower(orig) == strings.ToLower(ret) {
			ret = orig
		}
	}
	return ret
}

func StringForms(s string) (string, string) {
	return StringToSingular(s), StringToPlural(s)
}

func StringPlural(count int, s string) string {
	if count == 1 || count == -1 {
		return fmt.Sprintf("%d %s", count, StringToSingular(s))
	}
	return fmt.Sprintf("%d %s", count, StringToPlural(s))
}
