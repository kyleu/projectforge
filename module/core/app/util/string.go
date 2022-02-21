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

var acronyms = []string{"Id", "Api", "Json", "Html", "Xml"}

func StringToCamel(s string) string {
	return acr(strcase.ToCamel(s))
}

func StringToTitle(s string) string {
	ret := strings.Builder{}
	runes := []rune(StringToCamel(s))
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

func StringToLowerCamel(s string) string {
	return acr(strcase.ToLowerCamel(s))
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
