// Content managed by Project Forge, see [projectforge.md] for details.
package util

import (
	"encoding/xml"
	"fmt"
	"net/url"
	"strings"

	"github.com/pkg/errors"
	"golang.org/x/exp/slices"
)

type ValueMap map[string]any

func ValueMapFor(kvs ...any) ValueMap {
	ret := make(ValueMap, len(kvs)/2)
	ret.Add(kvs...)
	return ret
}

func (m ValueMap) KeysAndValues() ([]string, []any) {
	cols := make([]string, 0, len(m))
	vals := make([]any, 0, len(m))
	for k := range m {
		cols = append(cols, k)
	}
	slices.Sort(cols)
	for _, col := range cols {
		vals = append(vals, m[col])
	}
	return cols, vals
}

const selectedSuffix = "--selected"

func (m ValueMap) AsChanges() (ValueMap, error) {
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

	ret := make(ValueMap, len(keys))
	for _, k := range keys {
		ret[k] = vals[k]
	}
	return ret, nil
}

func (m ValueMap) Keys() []string {
	ret := make([]string, 0, len(m))
	for k := range m {
		ret = append(ret, k)
	}
	slices.Sort(ret)
	return ret
}

func (m ValueMap) Merge(args ...ValueMap) ValueMap {
	ret := make(ValueMap, len(m)+len(args))
	for k, v := range m {
		ret[k] = v
	}
	for _, arg := range args {
		for k, v := range arg {
			ret[k] = v
		}
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

func (m ValueMap) Clone() ValueMap {
	ret := make(ValueMap, len(m))
	for k, v := range m {
		ret[k] = v
	}
	return ret
}

func (m ValueMap) ToQueryString() string {
	params := url.Values{}
	for k, v := range m {
		params.Add(k, fmt.Sprint(v))
	}
	return params.Encode()
}

func (m ValueMap) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	tokens := []xml.Token{start}
	for key, value := range m {
		t := xml.StartElement{Name: xml.Name{Space: "", Local: key}}
		x, err := xml.Marshal(value)
		if err != nil {
			return err
		}
		tokens = append(tokens, t, xml.CharData(x), xml.EndElement{Name: t.Name})
	}
	tokens = append(tokens, xml.EndElement{Name: start.Name})
	for _, t := range tokens {
		err := e.EncodeToken(t)
		if err != nil {
			return err
		}
	}
	return e.Flush()
}

func (m ValueMap) Filter(keys []string) ValueMap {
	filteredMap := ValueMap{}
	for _, key := range keys {
		if data, ok := m[key]; ok {
			filteredMap[key] = data
		}
	}
	return filteredMap
}

func (m ValueMap) Overwrite(sourceMap ValueMap) ValueMap {
	destMap := m.Clone()
	for key, data := range sourceMap {
		destMap[key] = data
	}
	return destMap
}
