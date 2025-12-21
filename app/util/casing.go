package util

import (
	"strings"
	"unicode"

	"github.com/samber/lo"
)

func StringToProper(s string, extraAcronyms ...string) string {
	return acr(toProperCase(s, true), extraAcronyms...)
}

func StringToTitle(s string, extraAcronyms ...string) string {
	ret := strings.Builder{}
	runes := []rune(StringToProper(s, extraAcronyms...))
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

func StringToInitials(s string, extraAcronyms ...string) string {
	return StringJoin(lo.Map(strings.Split(StringToSnake(s, extraAcronyms...), "_"), func(s string, _ int) string {
		return strings.ToLower(s[0:1])
	}), "")
}

func StringToCamel(s string, extraAcronyms ...string) string {
	return acr(toProperCase(s, false), extraAcronyms...)
}

func StringToSnake(s string, extraAcronyms ...string) string {
	return acr(toDelimited(s, '_'), extraAcronyms...)
}

func StringToKebab(s string, extraAcronyms ...string) string {
	return acr(toDelimited(s, '-'), extraAcronyms...)
}

func toProperCase(s string, initCase bool) string {
	s = strings.TrimSpace(s)
	if s == "" {
		return s
	}
	if a, ok := acronymMap[s]; ok {
		s = a
	}

	n := strings.Builder{}
	n.Grow(len(s))
	capNext := initCase
	for i, v := range []byte(s) {
		vIsCap := v >= 'A' && v <= 'Z'
		vIsLow := v >= 'a' && v <= 'z'
		if capNext {
			if vIsLow {
				v += 'A'
				v -= 'a'
			}
		} else if i == 0 {
			if vIsCap {
				v += 'a'
				v -= 'A'
			}
		}
		if vIsCap || vIsLow {
			n.WriteByte(v)
			capNext = false
		} else if vIsNum := v >= '0' && v <= '9'; vIsNum {
			n.WriteByte(v)
			capNext = true
		} else {
			capNext = v == '_' || v == ' ' || v == '-' || v == '.'
		}
	}
	return n.String()
}

func toDelimited(s string, delimiter uint8) string {
	return toScreamingDelimited(s, delimiter, "", false)
}

//nolint:nestif,gocognit,gocyclo,cyclop
func toScreamingDelimited(s string, delimiter uint8, ignore string, screaming bool) string {
	s = strings.TrimSpace(s)
	n := strings.Builder{}
	n.Grow(len(s) + 2)
	for i, v := range []byte(s) {
		vIsCap := v >= 'A' && v <= 'Z'
		vIsLow := v >= 'a' && v <= 'z'
		if vIsLow && screaming {
			v += 'A'
			v -= 'a'
		} else if vIsCap && !screaming {
			v += 'a'
			v -= 'A'
		}

		if i+1 < len(s) {
			next := s[i+1]
			vIsNum := v >= '0' && v <= '9'
			nextIsCap := next >= 'A' && next <= 'Z'
			nextIsLow := next >= 'a' && next <= 'z'
			nextIsNum := next >= '0' && next <= '9'
			if (vIsCap && (nextIsLow || nextIsNum)) || (vIsLow && (nextIsCap || nextIsNum)) || (vIsNum && (nextIsCap || nextIsLow)) {
				prevIgnore := ignore != "" && i > 0 && strings.ContainsAny(string(s[i-1]), ignore)
				if !prevIgnore {
					if vIsCap && nextIsLow {
						if prevIsCap := i > 0 && s[i-1] >= 'A' && s[i-1] <= 'Z'; prevIsCap {
							n.WriteByte(delimiter)
						}
					}
					n.WriteByte(v)
					if vIsLow || vIsNum || nextIsNum {
						n.WriteByte(delimiter)
					}
					continue
				}
			}
		}

		if (v == ' ' || v == '_' || v == '-' || v == '.') && !strings.ContainsAny(string(v), ignore) {
			n.WriteByte(delimiter)
		} else {
			n.WriteByte(v)
		}
	}

	return n.String()
}
