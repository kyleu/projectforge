package util

import (
	"cmp"
	"encoding/xml"
	"fmt"
	"maps"
	"net/url"
	"reflect"
	"strings"

	"github.com/pkg/errors"
	"github.com/samber/lo"
)

const KeyMap = "map"

type ValueMap map[string]any

func ValueMapFor(kvs ...any) ValueMap {
	ret := make(ValueMap, len(kvs)/2)
	ret.Add(kvs...)
	return ret
}

func ValueMapFrom(m map[string]any) ValueMap {
	ret := make(ValueMap, len(m))
	for k, v := range m {
		ret[k] = v
	}
	return ret
}

func ValueMapFromAny(x any) (ValueMap, error) {
	switch t := x.(type) {
	case ValueMap:
		return t, nil
	case map[string]any:
		return ValueMapFrom(t), nil
	case ToMap:
		return t.ToMap(), nil
	default:
		rv := reflect.ValueOf(x)
		if rv.Kind() == reflect.Map {
			args := make([]any, 0, rv.Len()*2)
			iter := rv.MapRange()
			for iter.Next() {
				args = append(args, iter.Key().Interface(), iter.Value().Interface())
			}
			return ValueMapFor(args...), nil
		}
		return nil, errors.Errorf("unable to parse [%T] as ValueMap", x)
	}
}

func ValueMapFromAnyOK(x any) ValueMap {
	m, _ := ValueMapFromAny(x)
	return m
}

func (m ValueMap) Add(kvs ...any) {
	numWidth := -1
	pad := func() int {
		if numWidth == -1 {
			numWidth = len(fmt.Sprintf("%d", len(m)+len(kvs)/2))
		}
		return numWidth
	}
	for i := 0; i < len(kvs); i += 2 {
		v := kvs[i]
		var k string
		switch t := v.(type) {
		case string:
			k = t
		case fmt.Stringer:
			k = t.String()
		case int:
			k = fmt.Sprintf("%0*d", pad(), t)
		case int64:
			k = fmt.Sprintf("%0*d", pad(), t)
		case float64:
			k = fmt.Sprintf("%0*f", pad(), t)
		case bool:
			k = fmt.Sprintf("%t", t)
		default:
			k = fmt.Sprintf("error-invalid-type:%T", v)
		}
		m[k] = kvs[i+1]
	}
}

func (m ValueMap) With(k string, v any) ValueMap {
	x := m.Clone()
	x[k] = v
	return x
}

func (m ValueMap) Merge(args ...ValueMap) ValueMap {
	ret := m.Clone()
	lo.ForEach(args, func(arg ValueMap, _ int) {
		for k, v := range arg {
			ret[k] = v
		}
	})
	return ret
}

func (m ValueMap) Filter(keys []string) ValueMap {
	filteredMap := ValueMap{}
	lo.ForEach(keys, func(key string, _ int) {
		if data, ok := m[key]; ok {
			filteredMap[key] = data
		}
	})
	return filteredMap
}

func (m ValueMap) Overwrite(sourceMap ValueMap) ValueMap {
	destMap := m.Clone()
	for key, data := range sourceMap {
		destMap[key] = data
	}
	return destMap
}

func (m ValueMap) Length() int {
	return len(m)
}

func (m ValueMap) Empty() bool {
	return len(m) == 0
}

func (m ValueMap) HasKey(key string) bool {
	_, ok := m[key]
	return ok
}

func (m ValueMap) Keys() []string {
	return MapKeysSorted(m)
}

func (m ValueMap) KeysAndValues() ([]string, []any) {
	cols := m.Keys()
	vals := lo.Map(cols, func(col string, _ int) any {
		return m[col]
	})
	return cols, vals
}

func (m ValueMap) Clone() ValueMap {
	if m == nil {
		return nil
	}
	ret := make(ValueMap, len(m))
	for k, v := range m {
		ret[k] = v
	}
	return ret
}

func (m ValueMap) ReplaceEnvVars(logger Logger) ValueMap {
	ret := make(ValueMap, len(m))
	for k := range m {
		ret[k] = ReplaceEnvVars(m.GetStringOpt(k), logger)
	}
	return ret
}

func (m ValueMap) WithoutKeys(keys ...string) ValueMap {
	ret := m.Clone()
	for _, key := range keys {
		delete(ret, key)
	}
	return ret
}

func (m ValueMap) AsMap(simplify bool) map[string]any {
	if !simplify {
		return m
	}
	ret := make(map[string]any, len(m))
	for k, v := range m {
		ret[k], _ = simplifyValue(k, v)
	}
	return ret
}

func (m ValueMap) String() string {
	return ToJSONCompact(m)
}

func (m ValueMap) JSON() string {
	return ToJSON(m)
}

func (m ValueMap) ToMap() ValueMap {
	return m
}

func (m ValueMap) ToOrderedMap() *OrderedMap[any] {
	return OrderedMapFromMap(m, true)
}

func (m ValueMap) ToStringMap() map[string]string {
	return lo.MapValues(m, func(_ any, key string) string {
		return m.GetStringOpt(key)
	})
}

func (m ValueMap) ToQueryString() string {
	params := url.Values{}
	for k, v := range m {
		params.Add(k, fmt.Sprint(v))
	}
	return params.Encode()
}

func (m ValueMap) AsChanges() (ValueMap, error) {
	const selectedSuffix = "--selected"
	var keys []string
	vals := ValueMap{}

	for k, v := range m {
		if strings.HasSuffix(k, selectedSuffix) {
			key := strings.TrimSuffix(k, selectedSuffix)
			keys = append(keys, key)
		} else {
			curr, ok := vals[k]
			if ok {
				return nil, errors.Errorf("multiple values presented for [%s] (%T/%T)", k, curr, v)
			}
			vals[k] = v
		}
	}

	return lo.Associate(keys, func(k string) (string, any) {
		return k, vals[k]
	}), nil
}

func (m ValueMap) MarshalXML(e *xml.Encoder, _ xml.StartElement) error {
	err := e.EncodeToken(xml.StartElement{Name: xml.Name{Local: KeyMap}})
	if err != nil {
		return err
	}
	for key, value := range m {
		err = e.EncodeElement(value, xml.StartElement{Name: xml.Name{Local: key}})
		if err != nil {
			return err
		}
	}
	err = e.EncodeToken(xml.EndElement{Name: xml.Name{Local: KeyMap}})
	if err != nil {
		return err
	}
	return e.Flush()
}

func ValueMapGet[T any](m ValueMap, pth string) (T, error) {
	x, err := m.GetPath(pth, false)
	if err != nil {
		return DefaultValue[T](), err
	}
	ret, err := Cast[T](x)
	if err != nil {
		var df T
		return df, errors.Errorf("map value is of type [%T], expected [%T]", x, df)
	}
	return ret, nil
}

type ToMap interface {
	ToMap() ValueMap
}

type ToMaps interface {
	ToMaps() []ValueMap
}

func MapKeys[K comparable, V any](m map[K]V) []K {
	return lo.Keys(m)
}

func MapClone[K comparable, V any](m map[K]V) map[K]V {
	return maps.Clone(m)
}

func MapKeysSorted[K cmp.Ordered, V any](m map[K]V) []K {
	return ArraySorted(MapKeys(m))
}

func MapUpdateKeys[K cmp.Ordered, V any, NK cmp.Ordered](m map[K]V, fn func(k K, v V) NK) map[NK]V {
	ret := map[NK]V{}
	for k, v := range m {
		ret[fn(k, v)] = v
	}
	return ret
}

func MapUpdateValues[K comparable, V any, NV any](m map[K]V, fn func(k K, v V) (NV, error)) (map[K]NV, error) {
	ret := map[K]NV{}
	for k, v := range m {
		nv, err := fn(k, v)
		if err != nil {
			return nil, err
		}
		ret[k] = nv
	}
	return ret, nil
}
