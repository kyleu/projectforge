package util

import (
	"strings"
	"unicode"

	"github.com/iancoleman/strcase"
)

func StringSplit(s string, sep byte, cutc bool) (string, string) {
	i := strings.IndexByte(s, sep)
	if i < 0 {
		return s, ""
	}
	if cutc {
		return s[:i], s[i+1:]
	}
	return s[:i], s[i:]
}

func StringSplitLast(s string, sep byte, cutc bool) (string, string) {
	i := strings.LastIndexByte(s, sep)
	if i < 0 {
		return s, ""
	}
	if cutc {
		return s[:i], s[i+1:]
	}
	return s[:i], s[i:]
}

func StringSplitAndTrim(s string, delim string) []string {
	split := strings.Split(s, delim)
	ret := make([]string, 0, len(split))
	for _, x := range split {
		x = strings.TrimSpace(x)
		if len(x) > 0 {
			ret = append(ret, x)
		}
	}
	return ret
}

func StringPad(s string, size int) string {
	for len(s) < size {
		s += " "
	}
	return s
}

func StringToCamel(s string, extraAcronyms ...string) string {
	return acr(strcase.ToCamel(s), extraAcronyms...)
}

func StringToTitle(s string, extraAcronyms ...string) string {
	ret := strings.Builder{}
	runes := []rune(StringToCamel(s, extraAcronyms...))
	for idx, c := range runes {
		if idx > 0 && idx < len(runes)-1 && unicode.IsUpper(c) {
			if !unicode.IsUpper(runes[idx+1]) {
				ret.WriteRune(' ')
			} else if !unicode.IsUpper(runes[idx-1]) {
				ret.WriteRune(' ')
			}
		}
		ret.WriteRune(c)
	}
	return ret.String()
}

func StringToLowerCamel(s string, extraAcronyms ...string) string {
	return acr(strcase.ToLowerCamel(s), extraAcronyms...)
}

func StringToSnake(s string, extraAcronyms ...string) string {
	return acr(strcase.ToSnake(s), extraAcronyms...)
}

func StringTruncate(s string, max int) string {
	if max > len(s) {
		return s
	}
	return s[:strings.LastIndex(s[:max], " ")]
}

func StringRepeat(s string, n int) string {
	ret := strings.Builder{}
	for i := 0; i < n; i++ {
		ret.WriteString(s)
	}
	return ret.String()
}

var acronyms = func() []string {
	ret := []string{"Api", "Html", "Id", "Ip", "Json", "Xml", "Uri", "Url"}
	for _, x := range ret {
		strcase.ConfigureAcronym(strings.ToUpper(x), strings.ToLower(x))
	}
	return ret
}()

func acr(ret string, extraAcronyms ...string) string {
	proc := func(a string) {
		var lastIdx int
		for {
			i := strings.Index(ret[lastIdx:], a)
			if i == -1 {
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
	for _, a := range acronyms {
		proc(a)
	}
	for _, a := range extraAcronyms {
		proc(a)
	}
	return ret
}
