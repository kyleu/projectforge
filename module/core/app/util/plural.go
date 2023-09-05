package util

import (
	"fmt"
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
	if len(ret) < 3 {
		return ret
	}
	if ret[len(ret)-1] == 'S' {
		runes := []rune(ret)
		c2 := runes[len(runes)-2]
		c3 := runes[len(runes)-3]
		if unicode.IsUpper(c2) && unicode.IsUpper(c3) {
			runes[len(runes)-1] = 's'
			ret = string(runes)
		}
	}
	return ret
}

func StringToSingular(s string) string {
	plrlSvc()
	return plrl.Singular(s)
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
