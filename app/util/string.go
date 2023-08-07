// Content managed by Project Forge, see [projectforge.md] for details.
package util

import (
	"fmt"
	"hash/fnv"
	"path/filepath"
	"runtime"
	"strings"
	"unicode"

	"github.com/iancoleman/strcase"
	"github.com/pkg/errors"
	"github.com/samber/lo"
)

var winLB, saneLB = "\r\n", "\n"

var StringDefaultLinebreak = func() string {
	if runtime.GOOS == "windows" {
		return winLB
	}
	return saneLB
}()

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

func StringSplitLastOnly(s string, sep byte, cutc bool) string {
	i := strings.LastIndexByte(s, sep)
	if i < 0 {
		return s
	}
	if cutc {
		return s[i+1:]
	}
	return s[i:]
}

func StringSplitAndTrim(s string, delim string) []string {
	return lo.FilterMap(strings.Split(s, delim), func(x string, _ int) (string, bool) {
		x = strings.TrimSpace(x)
		return x, len(x) > 0
	})
}

func StringSplitPath(s string) (string, string) {
	i := strings.LastIndexByte(s, '/')
	if i == -1 {
		i = strings.LastIndexByte(s, '\\')
	}
	if i < 0 {
		return "", s
	}
	return s[:i], s[i+1:]
}

func StringSplitPathAndTrim(s string) []string {
	return StringSplitAndTrim(s, string(filepath.ListSeparator))
}

func StringDetectLinebreak(s string) string {
	if strings.Contains(s, winLB) {
		return winLB
	}
	return saneLB
}

func StringSplitLines(s string) []string {
	return strings.Split(s, StringDetectLinebreak(s))
}

func StringPad(s string, size int) string {
	return StringPadRight(s, size, ' ')
}

func StringPadRight(s string, size int, chr rune) string {
	sLen := len(s)
	if sLen >= size {
		return s
	}
	sb := strings.Builder{}
	sb.WriteString(s)
	for i := len(s); i < size; i++ {
		sb.WriteRune(chr)
	}
	return sb.String()
}

func StringPadLeft(s string, size int, chr rune) string {
	sLen := len(s)
	if sLen >= size {
		return s
	}
	sb := strings.Builder{}
	lo.Times(size-sLen, func(_ int) struct{} {
		sb.WriteRune(chr)
		return struct{}{}
	})
	sb.WriteString(s)
	return sb.String()
}

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

func StringTruncate(s string, max int) string {
	if max > len(s) {
		return s
	}
	return s[:strings.LastIndex(s[:max], " ")]
}

func StringRepeat(s string, n int) string {
	ret := strings.Builder{}
	lo.Times(n, func(_ int) struct{} {
		ret.WriteString(s)
		return struct{}{}
	})
	return ret.String()
}

func StringSubstringBetween(s string, l string, r string) string {
	li, ri := strings.Index(s, l), strings.Index(s, r)
	if li == -1 {
		return ""
	}
	lio := li + len(l)
	if ri == -1 {
		ri = len(s)
	}
	return s[lio:ri]
}

func StringReplaceBetween(s string, l string, r string, replacement string) (string, error) {
	li, ri := strings.Index(s, l), strings.Index(s, r)
	if li == -1 {
		return "", errors.Errorf("substring [%s] does not appear in the source", l)
	}
	lio := li + len(l)
	if ri == -1 {
		ri = len(s)
	}
	return s[:lio] + replacement + s[ri:], nil
}

func StringHash(s string) uint32 {
	h := fnv.New32a()
	_, _ = h.Write([]byte(s))
	return h.Sum32()
}

func CountryFlag(code string) string {
	if len(code) != 2 {
		return fmt.Sprintf("INVALID: %q", code)
	}
	code = strings.ToLower(code)
	const flagBaseIndex = '\U0001F1E6' - 'a'
	return string(rune(code[0])+flagBaseIndex) + string(rune(code[1])+flagBaseIndex)
}

var acronyms = func() []string {
	ret := []string{"Api", "Html", "Id", "Ip", "Json", "Xml", "Uri", "Url"}
	lo.ForEach(ret, func(x string, _ int) {
		strcase.ConfigureAcronym(strings.ToUpper(x), strings.ToLower(x))
	})
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
	lo.ForEach(acronyms, func(a string, _ int) {
		proc(a)
	})
	lo.ForEach(extraAcronyms, func(a string, _ int) {
		proc(a)
	})
	return ret
}
