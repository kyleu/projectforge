package util

import (
	"strings"
	"sync"
	"unicode"

	"github.com/samber/lo"
)

var (
	acronymMu  sync.RWMutex
	acronymMap = map[string]string{}
	acronyms   []string
)

func ConfigureAcronym(key, val string) {
	acronymMu.Lock()
	defer acronymMu.Unlock()
	acronymMap[key] = val
}

func InitAcronyms(extras ...string) error {
	acronymMu.Lock()
	defer acronymMu.Unlock()
	x := make([]string, 0, 12+len(extras))
	x = append(x, "Api", "Html", "Id", "Ip", "Json", "Md5", "Sha", "Sku", "Sql", "Xml", "Uri", "Url")
	x = append(x, lo.Map(extras, func(s string, _ int) string {
		return strings.ToUpper(s[:1]) + strings.ToLower(s[1:])
	})...)
	for _, a := range x {
		acronymMap[strings.ToUpper(a)] = strings.ToLower(a)
	}
	acronyms = x
	return nil
}

func Acronym(s string) string {
	acronymMu.Lock()
	defer acronymMu.Unlock()
	return acronymMap[s]
}

func acr(ret string, extraAcronyms ...string) string {
	acronymMu.RLock()
	activeAcronyms := append([]string(nil), acronyms...)
	acronymMu.RUnlock()
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
	for _, a := range activeAcronyms {
		proc(a)
	}
	for _, a := range extraAcronyms {
		proc(a)
	}
	return ret
}
