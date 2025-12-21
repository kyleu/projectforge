package util

import (
	"strconv"
	"strings"

	"github.com/google/uuid"
)

type RichString string

func RS(s string) RichString {
	return RichString(s)
}

func (r RichString) String() string {
	return string(r)
}

func (r RichString) Empty() bool {
	return len(r) == 0
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
	return RS(strings.ToLower(string(r)))
}

func (r RichString) ToUpper() RichString {
	return RS(strings.ToUpper(string(r)))
}

func (r RichString) TrimPrefix(prefix string) RichString {
	return RS(strings.TrimPrefix(string(r), prefix))
}

func (r RichString) TrimSuffix(suffix string) RichString {
	return RS(strings.TrimSuffix(string(r), suffix))
}

func (r RichString) TrimSpace() RichString {
	return RS(strings.TrimSpace(string(r)))
}

func (r RichString) Split(sep byte, cutc bool) (RichString, RichString) {
	left, right := StringSplit(string(r), sep, cutc)
	return RS(left), RS(right)
}

func (r RichString) SplitLast(sep byte, cutc bool) (RichString, RichString) {
	left, right := StringSplitLast(string(r), sep, cutc)
	return RS(left), RS(right)
}

func (r RichString) SplitLastOnly(sep byte, cutc bool) RichString {
	return RS(StringSplitLastOnly(string(r), sep, cutc))
}

func (r RichString) SplitAndTrim(delim string) RichStrings {
	return RStrings(StringSplitAndTrim(string(r), delim)...)
}

func (r RichString) SplitPath() (string, string) {
	return StringSplitPath(string(r))
}

func (r RichString) SplitPathAndTrim() RichStrings {
	return RStrings(StringSplitPathAndTrim(string(r))...)
}

func (r RichString) DetectLinebreak() string {
	return StringDetectLinebreak(string(r))
}

func (r RichString) SplitLines() RichStrings {
	return RStrings(StringSplitLines(string(r))...)
}

func (r RichString) SplitLinesIndented(indent int, indentFirstLine bool, includeEmptyLines bool) RichStrings {
	return RStrings(StringSplitLinesIndented(string(r), indent, indentFirstLine, includeEmptyLines)...)
}

func (r RichString) Pad(size int) RichString {
	return RS(StringPadRight(string(r), size, ' '))
}

func (r RichString) PadRight(size int, chr rune) RichString {
	return RS(StringPadRight(string(r), size, chr))
}

func (r RichString) PadLeft(size int, chr rune) RichString {
	return RS(StringPadLeft(string(r), size, chr))
}

func (r RichString) Truncate(mx int) RichString {
	return RS(StringTruncate(string(r), mx))
}

func (r RichString) Repeat(n int) RichString {
	return RS(StringRepeat(string(r), n))
}

func (r RichString) SubstringBetween(left string, right string) RichString {
	return RS(StringSubstringBetween(string(r), left, right))
}

func (r RichString) ReplaceBetween(left string, right string, replacement string) (RichString, error) {
	ret, err := StringReplaceBetween(string(r), left, right, replacement)
	return RS(ret), err
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
	return r.With(RStrings(elems...)...)
}

func (r RichString) Join(elems RichStrings, delim string) RichString {
	return RS(StringJoin(r.With(elems...).Strings(), delim))
}

func (r RichString) JoinStrings(elems []string, delim string) RichString {
	return RS(StringJoin(r.With(RStrings(elems...)...).Strings(), delim))
}

func (r RichString) Path(elems ...string) RichString {
	return RS(StringPath(r.WithStrings(elems...).Strings()...))
}

func (r RichString) FilePath(elems ...string) RichString {
	return RS(StringFilePath(r.WithStrings(elems...).Strings()...))
}

func (r RichString) OrDefault(dflt string) RichString {
	if r.Empty() {
		return RS(dflt)
	}
	return r
}

func (r RichString) Append(strs ...string) RichString {
	ret := r.String()
	for _, s := range strs {
		ret += s
	}
	return RS(ret)
}

type RichStrings []RichString

func RStrings(strs ...string) RichStrings {
	ret := make(RichStrings, 0, len(strs))
	for _, s := range strs {
		ret = append(ret, RS(s))
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
	return RS(StringJoin(r.Strings(), delim))
}
