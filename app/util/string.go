package util

import (
	"fmt"
	"path/filepath"
	"reflect"
	"runtime"
	"strings"

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
		return x, x != ""
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

func StringSplitLinesIndented(s string, indent int, indentFirstLine bool, includeEmptyLines bool) []string {
	ind := StringRepeat("  ", indent)
	lines := StringSplitLines(s)
	ret := NewStringSlice(make([]string, 0, len(lines)))
	for idx, line := range lines {
		if (!includeEmptyLines) && strings.TrimSpace(line) == "" {
			continue
		}
		if indentFirstLine || idx > 0 {
			line = ind + line
		}
		ret.Push(line)
	}
	return ret.Slice
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
		return EmptyStruct
	})
	sb.WriteString(s)
	return sb.String()
}

func StringTruncate(s string, mx int) string {
	if mx > len(s) {
		return s
	}
	return s[:strings.LastIndex(s[:mx], " ")]
}

func StringRepeat(s string, n int) string {
	ret := strings.Builder{}
	lo.Times(n, func(_ int) struct{} {
		ret.WriteString(s)
		return EmptyStruct
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

func CountryFlag(code string) string {
	if len(code) != 2 {
		return fmt.Sprintf("INVALID: %q", code)
	}
	code = strings.ToLower(code)
	const flagBaseIndex = '\U0001F1E6' - 'a'
	return string(rune(code[0])+flagBaseIndex) + string(rune(code[1])+flagBaseIndex)
}

var filenameReplacer = strings.NewReplacer("/", "-", "\\", "-", "?", "-", "%", "-", "*", "-", ":", "-", "|", "-", "\"", "-", "<", "-", ">", "-")

func Filename(s string) string {
	return filenameReplacer.Replace(s)
}

func StringNullable(s fmt.Stringer) string {
	if s == nil || reflect.ValueOf(s).IsNil() {
		return ""
	}
	return s.String()
}
