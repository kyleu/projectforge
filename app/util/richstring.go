package util

import (
	"strconv"
	"strings"

	"github.com/google/uuid"
)

type Str string

func (s Str) String() string {
	return string(s)
}

func (s Str) Empty() bool {
	return len(s) == 0
}

func (s Str) Equal(x string) bool {
	return string(s) == x
}

func (s Str) EqualFold(x string) bool {
	return strings.EqualFold(string(s), x)
}

func (s Str) Length() int {
	return len(s)
}

func (s Str) Index(substring string) int {
	return strings.Index(s.String(), substring)
}

func (s Str) Contains(x string) bool {
	return strings.Contains(string(s), x)
}

func (s Str) ContainsAny(x string) bool {
	return strings.ContainsAny(string(s), x)
}

func (s Str) ContainsRune(chr rune) bool {
	return strings.ContainsRune(string(s), chr)
}

func (s Str) HasPrefix(prefix string) bool {
	return strings.HasPrefix(string(s), prefix)
}

func (s Str) HasSuffix(suffix string) bool {
	return strings.HasSuffix(string(s), suffix)
}

func (s Str) ToLower() Str {
	return Str(strings.ToLower(string(s)))
}

func (s Str) ToUpper() Str {
	return Str(strings.ToUpper(string(s)))
}

func (s Str) TrimPrefix(prefix string) Str {
	return Str(strings.TrimPrefix(string(s), prefix))
}

func (s Str) TrimSuffix(suffix string) Str {
	return Str(strings.TrimSuffix(string(s), suffix))
}

func (s Str) TrimSpace() Str {
	return Str(strings.TrimSpace(string(s)))
}

func (s Str) Cut(sep byte, cutc bool) (Str, Str) {
	left, right := StringCut(string(s), sep, cutc)
	return Str(left), Str(right)
}

func (s Str) CutLast(sep byte, cutc bool) (Str, Str) {
	left, right := StringCutLast(string(s), sep, cutc)
	return Str(left), Str(right)
}

func (s Str) Split(delim string) Strings {
	return Strs(strings.Split(s.String(), delim)...)
}

func (s Str) SplitLastOnly(sep byte, cutc bool) Str {
	return Str(StringSplitLastOnly(string(s), sep, cutc))
}

func (s Str) SplitAndTrim(delim string) Strings {
	return Strs(StringSplitAndTrim(string(s), delim)...)
}

func (s Str) SplitPath() (string, string) {
	return StringSplitPath(string(s))
}

func (s Str) SplitPathAndTrim() Strings {
	return Strs(StringSplitPathAndTrim(string(s))...)
}

func (s Str) DetectLinebreak() string {
	return StringDetectLinebreak(string(s))
}

func (s Str) SplitLines() Strings {
	return Strs(StringSplitLines(string(s))...)
}

func (s Str) SplitLinesIndented(indent int, indentFirstLine bool, includeEmptyLines bool) Strings {
	return Strs(StringSplitLinesIndented(string(s), indent, indentFirstLine, includeEmptyLines)...)
}

func (s Str) Pad(size int) Str {
	return Str(StringPadRight(string(s), size, ' '))
}

func (s Str) PadRight(size int, chr rune) Str {
	return Str(StringPadRight(string(s), size, chr))
}

func (s Str) PadLeft(size int, chr rune) Str {
	return Str(StringPadLeft(string(s), size, chr))
}

func (s Str) Truncate(mx int) Str {
	return Str(StringTruncate(string(s), mx))
}

func (s Str) Repeat(n int) Str {
	return Str(StringRepeat(string(s), n))
}

func (s Str) Substring(start int, end int) Str {
	return Str(string(s)[start:end])
}

func (s Str) SubstringBetween(left string, right string) Str {
	return Str(StringSubstringBetween(string(s), left, right))
}

func (s Str) Replace(old string, nw string, n int) Str {
	return Str(strings.Replace(string(s), old, nw, n))
}

func (s Str) ReplaceAll(old string, nw string) Str {
	return Str(strings.ReplaceAll(string(s), old, nw))
}

func (s Str) ReplaceBetween(left string, right string, replacement string) (Str, error) {
	ret, err := StringReplaceBetween(string(s), left, right, replacement)
	return Str(ret), err
}

func (s Str) ParseUUID() *uuid.UUID {
	return UUIDFromString(s.String())
}

func (s Str) ParseInt() (int, bool) {
	ret, err := strconv.ParseInt(s.String(), 10, 32)
	return int(ret), err == nil
}

func (s Str) WithPrefix(prefixes ...Str) Str {
	var ret Strings = append(ArrayCopy(prefixes), s)
	return ret.Join("")
}

func (s Str) WithSuffix(suffixes ...Str) Str {
	return s.With(suffixes...).Join("")
}

func (s Str) With(elems ...Str) Strings {
	return append(Strings{s}, elems...)
}

func (s Str) WithStrings(elems ...string) Strings {
	return s.With(Strs(elems...)...)
}

func (s Str) Join(elems Strings, delim string) Str {
	return Str(StringJoin(s.With(elems...).Strings(), delim))
}

func (s Str) JoinStrings(elems []string, delim string) Str {
	return Str(StringJoin(s.With(Strs(elems...)...).Strings(), delim))
}

func (s Str) Path(elems ...string) Str {
	return Str(StringPath(s.WithStrings(elems...).Strings()...))
}

func (s Str) FilePath(elems ...string) Str {
	return Str(StringFilePath(s.WithStrings(elems...).Strings()...))
}

func (s Str) ToProper() Str {
	return Str(StringToProper(string(s)))
}

func (s Str) ToCamel() Str {
	return Str(StringToCamel(string(s)))
}

func (s Str) ToSnake() Str {
	return Str(StringToSnake(string(s)))
}

func (s Str) ToKebab() Str {
	return Str(StringToKebab(string(s)))
}

func (s Str) OrDefault(dflt string) Str {
	if s.Empty() {
		return Str(dflt)
	}
	return s
}

func (s Str) Append(strs ...string) Str {
	ret := s.String()
	for _, s := range strs {
		ret += s
	}
	return Str(ret)
}

type Strings []Str

func Strs(strs ...string) Strings {
	ret := make(Strings, 0, len(strs))
	for _, s := range strs {
		ret = append(ret, Str(s))
	}
	return ret
}

func (s Strings) Strings() []string {
	ret := make([]string, 0, len(s))
	for _, x := range s {
		ret = append(ret, x.String())
	}
	return ret
}

func (s Strings) Empty() bool {
	return len(s) == 0
}

func (s Strings) Join(delim string) Str {
	return Str(StringJoin(s.Strings(), delim))
}
