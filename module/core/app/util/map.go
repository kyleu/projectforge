package util

import (
	"encoding/xml"
	"fmt"
	"net/url"
	"strings"

	"github.com/pkg/errors"
	"github.com/samber/lo"
)

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

func (m ValueMap) Add(kvs ...any) {
	for i := 0; i < len(kvs); i += 2 {
		k, ok := kvs[i].(string)
		if !ok {
			k = fmt.Sprintf("error-invalid-type:%T", kvs[i])
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

func (m ValueMap) HasKey(key string) bool {
	_, ok := m[key]
	return ok
}

func (m ValueMap) Keys() []string {
	return ArraySorted(lo.Keys(m))
}

func (m ValueMap) KeysAndValues() ([]string, []any) {
	cols := ArraySorted(lo.Keys(m))
	vals := lo.Map(cols, func(col string, _ int) any {
		return m[col]
	})
	return cols, vals
}

func (m ValueMap) Clone() ValueMap {
	ret := make(ValueMap, len(m))
	for k, v := range m {
		ret[k] = v
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
		ret[k] = simplifyValue(k, v)
	}
	return ret
}

func (m ValueMap) String() string {
	return ToJSONCompact(m)
}

func (m ValueMap) JSON() string {
	return ToJSON(m)
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
	err := e.EncodeToken(xml.StartElement{Name: xml.Name{Local: "map"}})
	if err != nil {
		return err
	}
	for key, value := range m {
		err = e.EncodeElement(value, xml.StartElement{Name: xml.Name{Local: key}})
		if err != nil {
			return err
		}
	}
	err = e.EncodeToken(xml.EndElement{Name: xml.Name{Local: "map"}})
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
	ret, ok := x.(T)
	if !ok {
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
