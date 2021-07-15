package util

import (
	"strings"
	"unicode"

	"github.com/iancoleman/strcase"
)

func SplitString(s string, sep byte, cutc bool) (string, string) {
	i := strings.IndexByte(s, sep)
	if i < 0 {
		return s, ""
	}
	if cutc {
		return s[:i], s[i+1:]
	}
	return s[:i], s[i:]
}

func SplitStringLast(s string, sep byte, cutc bool) (string, string) {
	i := strings.LastIndexByte(s, sep)
	if i < 0 {
		return s, ""
	}
	if cutc {
		return s[:i], s[i+1:]
	}
	return s[:i], s[i:]
}

func SplitAndTrim(s string, delim string) []string {
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

var acronyms = []string{"Id"}

func ToCamel(s string) string {
	return acr(strcase.ToCamel(s))
}

func ToTitle(s string) string {
	ret := strings.Builder{}
	runes := []rune(ToCamel(s))
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

func ToLowerCamel(s string) string {
	return acr(strcase.ToLowerCamel(s))
}

func acr(ret string) string {
	for _, a := range acronyms {
		for {
			i := strings.Index(ret, a)
			if i == -1 {
				break
			}
			ret = ret[:i] + strings.ToUpper(a) + ret[i+len(a):]
		}
	}
	return ret
}

func OxfordComma(names []string, separator string) string {
	ret := ""
	for idx, name := range names {
		if idx > 0 {
			if idx == (len(names) - 1) {
				if idx > 1 {
					ret += ","
				}
				ret += " " + separator + " "
			} else {
				ret += ", "
			}
		}
		ret += name
	}
	return ret
}
