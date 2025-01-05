package util

import (
	"encoding/csv"
	"strconv"
	"strings"

	"github.com/pkg/errors"
)

func (m ValueMap) GetPath(path string, allowMissing bool) (any, error) {
	match, ok := m[path]
	if ok {
		return match, nil
	}
	r := csv.NewReader(strings.NewReader(path)) // to support quoted strings like files."readme.txt".size
	r.Comma = '.'
	fields, err := r.Read()
	if err != nil {
		return nil, err
	}
	ret, err := getPath(m, allowMissing, fields...)
	if err != nil {
		return nil, errors.Wrapf(err, "invalid path [%s]", path)
	}
	return ret, nil
}

func getPath(i any, allowMissing bool, path ...string) (any, error) {
	if len(path) == 0 {
		return i, nil
	}
	k := path[0]
	switch t := i.(type) {
	case ValueMap:
		ret, ok := t[k]
		if !ok {
			if allowMissing {
				return nil, nil
			}
			return nil, errors.Errorf("map does not have key [%s] among candidates [%s]", k, strings.Join(t.Keys(), ", "))
		}
		return getPath(ret, allowMissing, path[1:]...)
	case map[string]any:
		ret, ok := t[k]
		if !ok {
			if allowMissing {
				return nil, nil
			}
			return nil, errors.Errorf("map does not have key [%s] among candidates [%s]", k, strings.Join(MapKeys(t), ", "))
		}
		return getPath(ret, allowMissing, path[1:]...)
	case []any:
		i, err := strconv.ParseInt(k, 10, 32)
		if err != nil {
			return nil, errors.Errorf("path [%s] refers to an slice, but can't be parsed as an index", k)
		}
		var ret any
		if len(t) > int(i) {
			ret = t[int(i)]
		}
		return getPath(ret, allowMissing, path[1:]...)
	default:
		if allowMissing {
			return nil, nil
		}
		return nil, errors.Errorf("unhandled type [%T] for path [%s]", k, strings.Join(path, "."))
	}
}

func (m ValueMap) SetPath(path string, val any) error {
	if !strings.Contains(path, ".") {
		m[path] = val
		return nil
	}
	r := csv.NewReader(strings.NewReader(path))
	r.Comma = '.'
	fields, err := r.Read()
	if err != nil {
		return err
	}
	return setPath(m, fields, val)
}

func setPath(i any, path []string, val any) error {
	work := i
	for idx, p := range path {
		if idx == len(path)-1 {
			switch t := work.(type) {
			case ValueMap:
				t[p] = val
			case map[string]any:
				t[p] = val
			default:
				return errors.Errorf("unhandled [%T]", t)
			}
		} else {
			switch t := work.(type) {
			case ValueMap:
				t[p] = map[string]any{}
				work = t[p]
			case map[string]any:
				t[p] = map[string]any{}
				work = t[p]
			default:
				return errors.Errorf("unhandled [%T]", t)
			}
		}
	}
	return nil
}
