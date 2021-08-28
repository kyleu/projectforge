package util

import (
	"encoding/csv"
	"fmt"
	"net/url"
	"sort"
	"strconv"
	"strings"

	"github.com/pkg/errors"
)

type ValueMap map[string]interface{}

func ValueMapFor(kvs ...interface{}) ValueMap {
	ret := make(ValueMap, len(kvs)/2)
	ret.Add(kvs...)
	return ret
}

func (m ValueMap) KeysAndValues() ([]string, []interface{}) {
	cols := make([]string, 0, len(m))
	vals := make([]interface{}, 0, len(m))
	for k := range m {
		cols = append(cols, k)
	}
	sort.Strings(cols)
	for _, col := range cols {
		vals = append(vals, m[col])
	}
	return cols, vals
}

func (m ValueMap) GetRequired(k string) (interface{}, error) {
	v, ok := m[k]
	if !ok {
		msg := "no value [%s] among candidates [%s]"
		return nil, errors.Errorf(msg, k, OxfordComma(m.Keys(), "and"))
	}
	return v, nil
}

func (m ValueMap) GetBool(k string) (bool, error) {
	v, err := m.GetRequired(k)
	if err != nil {
		return false, err
	}

	var ret bool
	switch t := v.(type) {
	case bool:
		ret = t
	case string:
		ret = t == "true"
	case nil:
		ret = false
	default:
		return false, errors.Errorf("expected boolean or string, encountered %T", t)
	}
	return ret, nil
}

func (m ValueMap) GetInteger(k string, allowEmpty bool) (int, error) {
	v, err := m.GetRequired(k)
	if err != nil {
		return 0, err
	}

	switch t := v.(type) {
	case int:
		return t, nil
	case int32:
		return int(t), nil
	case int64:
		return int(t), nil
	case nil:
		if allowEmpty {
			return 0, nil
		}
		return 0, errors.New(k + " is nil, not integer")
	default:
		return 0, errors.Errorf("expected integer, encountered %T", t)
	}
}

func (m ValueMap) GetString(k string, allowEmpty bool) (string, error) {
	v, err := m.GetRequired(k)
	if err != nil {
		return "", err
	}

	var ret string
	switch t := v.(type) {
	case []string:
		ret = strings.Join(t, "|")
	case string:
		ret = t
	case nil:
		ret = ""
	default:
		return "", errors.Errorf("expected string or array of strings, encountered %T", t)
	}
	if !allowEmpty && ret == "" {
		return "", errors.Errorf("field [%s] may not be empty", k)
	}
	return ret, nil
}

func (m ValueMap) GetStringOpt(k string) string {
	ret, _ := m.GetString(k, true)
	return ret
}

func (m ValueMap) GetStringArray(k string, allowMissing bool) ([]string, error) {
	v, err := m.GetRequired(k)
	if err != nil {
		if allowMissing {
			return nil, nil
		}
		return nil, err
	}

	switch t := v.(type) {
	case []string:
		return t, nil
	case string:
		return []string{t}, nil
	default:
		return nil, errors.Errorf("expected array of strings, encountered %T", t)
	}
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
	sort.Strings(ret)
	return ret
}

func (m ValueMap) Merge(args ValueMap) ValueMap {
	ret := make(ValueMap, len(m)+len(args))
	for k, v := range m {
		ret[k] = v
	}
	for k, v := range args {
		ret[k] = v
	}
	return ret
}

func (m ValueMap) Add(kvs ...interface{}) {
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

func (m ValueMap) GetPath(path string) interface{} {
	r := csv.NewReader(strings.NewReader(path))
	r.Comma = '.'
	fields, err := r.Read()
	if err != nil {
		return err
	}
	return getPath(m, fields)
}

func getPath(i interface{}, path []string) interface{} {
	if len(path) == 0 {
		return i
	}
	k := path[0]
	switch t := i.(type) {
	case ValueMap:
		ret, ok := t[k]
		if !ok {
			return nil
		}
		return getPath(ret, path[1:])
	case map[string]interface{}:
		ret, ok := t[k]
		if !ok {
			return nil
		}
		return getPath(ret, path[1:])
	case []interface{}:
		i, err := strconv.Atoi(k)
		if err != nil {
			return nil
		}
		var ret interface{}
		if len(t) > i {
			ret = t[i]
		}
		return getPath(ret, path[1:])
	default:
		return nil
	}
}

func (m ValueMap) SetPath(path string, val interface{}) interface{} {
	r := csv.NewReader(strings.NewReader(path))
	r.Comma = '.'
	fields, err := r.Read()
	if err != nil {
		return err
	}
	return setPath(m, fields, val)
}

func (m ValueMap) Unset(s string) {
	delete(m, s)
}

func setPath(i interface{}, path []string, val interface{}) error {
	work := i
	for idx, p := range path {
		if idx == len(path)-1 {
			switch t := work.(type) {
			case ValueMap:
				t[p] = val
			case map[string]interface{}:
				t[p] = val
			default:
				return errors.New(fmt.Sprintf("unhandled [%T]", t))
			}
		} else {
			switch t := work.(type) {
			case ValueMap:
				t[p] = map[string]interface{}{}
				work = t[p]
			case map[string]interface{}:
				t[p] = map[string]interface{}{}
				work = t[p]
			default:
				return errors.New(fmt.Sprintf("unhandled [%T]", t))
			}
		}
	}
	return nil
}
