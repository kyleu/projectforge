package util

import (
	"strconv"
	"strings"

	"github.com/google/uuid"
)

type RichString string

func Str(s string) RichString {
	return RichString(s)
}

func (r RichString) String() string {
	return string(r)
}

func (r RichString) Empty() bool {
	return len(r) == 0
}

func (r RichString) Equal(s string) bool {
	return string(r) == s
}

func (r RichString) EqualFold(s string) bool {
	return strings.EqualFold(string(r), s)
}

func (r RichString) Length() int {
	return len(r)
}

func (r RichString) Index(substring string) int {
	return strings.Index(r.String(), substring)
}

func (r RichString) Contains(s string) bool {
	return strings.Contains(string(r), s)
}

func (r RichString) ContainsAny(s string) bool {
	return strings.ContainsAny(string(r), s)
}

func (r RichString) ContainsRune(chr rune) bool {
	return strings.ContainsRune(string(r), chr)
}

func (r RichString) HasPrefix(prefix string) bool {
	return strings.HasPrefix(string(r), prefix)
}

func (r RichString) HasSuffix(suffix string) bool {
	return strings.HasSuffix(string(r), suffix)
}

func (r RichString) ToLower() RichString {
	return Str(strings.ToLower(string(r)))
}

func (r RichString) ToUpper() RichString {
	return Str(strings.ToUpper(string(r)))
}

func (r RichString) TrimPrefix(prefix string) RichString {
	return Str(strings.TrimPrefix(string(r), prefix))
}

func (r RichString) TrimSuffix(suffix string) RichString {
	return Str(strings.TrimSuffix(string(r), suffix))
}

func (r RichString) TrimSpace() RichString {
	return Str(strings.TrimSpace(string(r)))
}

func (r RichString) Cut(sep byte, cutc bool) (RichString, RichString) {
	left, right := StringCut(string(r), sep, cutc)
	return Str(left), Str(right)
}

func (r RichString) CutLast(sep byte, cutc bool) (RichString, RichString) {
	left, right := StringCutLast(string(r), sep, cutc)
	return Str(left), Str(right)
}

func (r RichString) Split(delim string) RichStrings {
	return Strs(strings.Split(r.String(), delim)...)
}

func (r RichString) SplitLastOnly(sep byte, cutc bool) RichString {
	return Str(StringSplitLastOnly(string(r), sep, cutc))
}

func (r RichString) SplitAndTrim(delim string) RichStrings {
	return Strs(StringSplitAndTrim(string(r), delim)...)
}

func (r RichString) SplitPath() (string, string) {
	return StringSplitPath(string(r))
}

func (r RichString) SplitPathAndTrim() RichStrings {
	return Strs(StringSplitPathAndTrim(string(r))...)
}

func (r RichString) DetectLinebreak() string {
	return StringDetectLinebreak(string(r))
}

func (r RichString) SplitLines() RichStrings {
	return Strs(StringSplitLines(string(r))...)
}

func (r RichString) SplitLinesIndented(indent int, indentFirstLine bool, includeEmptyLines bool) RichStrings {
	return Strs(StringSplitLinesIndented(string(r), indent, indentFirstLine, includeEmptyLines)...)
}

func (r RichString) Pad(size int) RichString {
	return Str(StringPadRight(string(r), size, ' '))
}

func (r RichString) PadRight(size int, chr rune) RichString {
	return Str(StringPadRight(string(r), size, chr))
}

func (r RichString) PadLeft(size int, chr rune) RichString {
	return Str(StringPadLeft(string(r), size, chr))
}

func (r RichString) Truncate(mx int) RichString {
	return Str(StringTruncate(string(r), mx))
}

func (r RichString) Repeat(n int) RichString {
	return Str(StringRepeat(string(r), n))
}

func (r RichString) Substring(start int, end int) RichString {
	return Str(string(r)[start:end])
}

func (r RichString) SubstringBetween(left string, right string) RichString {
	return Str(StringSubstringBetween(string(r), left, right))
}

func (r RichString) Replace(old string, new string, n int) RichString {
	return Str(strings.Replace(string(r), old, new, n))
}

func (r RichString) ReplaceAll(old string, new string) RichString {
	return Str(strings.ReplaceAll(string(r), old, new))
}

func (r RichString) ReplaceBetween(left string, right string, replacement string) (RichString, error) {
	ret, err := StringReplaceBetween(string(r), left, right, replacement)
	return Str(ret), err
}

func (r RichString) ParseUUID() *uuid.UUID {
	return UUIDFromString(r.String())
}

func (r RichString) ParseInt() (int, bool) {
	ret, err := strconv.ParseInt(r.String(), 10, 32)
	return int(ret), err == nil
}

func (r RichString) WithPrefix(prefixes ...RichString) RichString {
	var ret RichStrings = append(ArrayCopy(prefixes), r)
	return ret.Join("")
}

func (r RichString) WithSuffix(suffixes ...RichString) RichString {
	var ret RichStrings = append(RichStrings{r}, suffixes...)
	return ret.Join("")
}

func (r RichString) With(elems ...RichString) RichStrings {
	return append(RichStrings{r}, elems...)
}

func (r RichString) WithStrings(elems ...string) RichStrings {
	return r.With(Strs(elems...)...)
}

func (r RichString) Join(elems RichStrings, delim string) RichString {
	return Str(StringJoin(r.With(elems...).Strings(), delim))
}

func (r RichString) JoinStrings(elems []string, delim string) RichString {
	return Str(StringJoin(r.With(Strs(elems...)...).Strings(), delim))
}

func (r RichString) Path(elems ...string) RichString {
	return Str(StringPath(r.WithStrings(elems...).Strings()...))
}

func (r RichString) FilePath(elems ...string) RichString {
	return Str(StringFilePath(r.WithStrings(elems...).Strings()...))
}

func (r RichString) ToProper() RichString {
	return Str(StringToProper(string(r)))
}

func (r RichString) ToCamel() RichString {
	return Str(StringToCamel(string(r)))
}

func (r RichString) ToSnake() RichString {
	return Str(StringToSnake(string(r)))
}

func (r RichString) ToKebab() RichString {
	return Str(StringToKebab(string(r)))
}

func (r RichString) OrDefault(dflt string) RichString {
	if r.Empty() {
		return Str(dflt)
	}
	return r
}

func (r RichString) Append(strs ...string) RichString {
	ret := r.String()
	for _, s := range strs {
		ret += s
	}
	return Str(ret)
}

type RichStrings []RichString

func Strs(strs ...string) RichStrings {
	ret := make(RichStrings, 0, len(strs))
	for _, s := range strs {
		ret = append(ret, Str(s))
	}
	return ret
}

func (r RichStrings) Strings() []string {
	ret := make([]string, 0, len(r))
	for _, s := range r {
		ret = append(ret, s.String())
	}
	return ret
}

func (r RichStrings) Empty() bool {
	return len(r) == 0
}

func (r RichStrings) Join(delim string) RichString {
	return Str(StringJoin(r.Strings(), delim))
}
