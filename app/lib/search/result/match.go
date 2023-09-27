// Package result - Content managed by Project Forge, see [projectforge.md] for details.
package result

import (
	"cmp"
	"fmt"
	"reflect"
	"slices"
	"strings"

	"github.com/samber/lo"
)

type Match struct {
	Key   string `json:"k"`
	Value string `json:"v"`
}

func (m *Match) ValueSplit(q string) []string {
	ql := strings.ToLower(q)
	vl := strings.ToLower(m.Value)
	cut := m.Value
	idx := strings.Index(vl, ql)
	if idx == -1 {
		return []string{cut}
	}
	var ret []string
	for idx > -1 {
		if idx > 0 {
			ret = append(ret, cut[:idx])
		}
		ret = append(ret, cut[idx:idx+len(ql)])

		cut = cut[idx+len(ql):]
		vl = vl[idx+len(ql):]

		idx = strings.Index(vl, ql)
	}
	if len(cut) > 0 {
		ret = append(ret, cut)
	}
	return ret
}

type Matches []*Match

func (m Matches) Sort() {
	slices.SortFunc(m, func(l *Match, r *Match) int {
		return cmp.Compare(strings.ToLower(l.Key), strings.ToLower(r.Key))
	})
}

func MatchesFor(key string, x any, q string) Matches {
	q = strings.ToLower(q)
	v := reflect.ValueOf(x)
	if v.Kind() == reflect.Ptr {
		if v.IsNil() {
			return nil
		}
		v = reflect.Indirect(v)
	}

	appendKey := func(s string) string {
		if key == "" {
			return s
		}
		return key + "." + s
	}
	maybe := func(cond bool, v string) Matches {
		if cond {
			return Matches{{Key: key, Value: v}}
		}
		return nil
	}

	switch v.Kind() {
	case reflect.Bool:
		return nil
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		i := fmt.Sprint(v.Int())
		return maybe(strings.Contains(i, q), i)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		i := fmt.Sprint(v.Uint())
		return maybe(strings.Contains(i, q), i)
	case reflect.Float32, reflect.Float64:
		f := fmt.Sprint(v.Float())
		return maybe(strings.Contains(f, q), f)
	case reflect.Map:
		var ret Matches
		x := v.MapRange()
		for x.Next() {
			ret = append(ret, MatchesFor(appendKey(x.Key().String()), x.Value().Interface(), q)...)
		}
		return ret
	case reflect.Array, reflect.Slice:
		var ret Matches
		for idx := 0; idx < v.Len(); idx++ {
			ret = append(ret, MatchesFor(appendKey(fmt.Sprint(idx)), v.Index(idx), q)...)
		}
		return ret
	case reflect.String:
		s := v.String()
		return maybe(strings.Contains(strings.ToLower(s), q), s)
	case reflect.Struct:
		var ret Matches
		lo.Times(v.NumField(), func(i int) struct{} {
			if f := v.Field(i); f.CanSet() {
				n := v.Type().Field(i).Name
				if m := MatchesFor(appendKey(n), v.Field(i).Interface(), q); m != nil {
					ret = append(ret, m...)
				}
			}
			return struct{}{}
		})
		return ret
	default:
		return Matches{{Key: key, Value: "error: " + v.Kind().String()}}
	}
}
