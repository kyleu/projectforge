package util

import (
	"encoding/csv"
	"strconv"
	"strings"

	"github.com/pkg/errors"
)

func (m ValueMap) GetPath(path string) interface{} {
	r := csv.NewReader(strings.NewReader(path))
	r.Comma = '.'
	fields, err := r.Read()
	if err != nil {
		return err
	}
	return getPath(m, fields...)
}

func getPath(i interface{}, path ...string) interface{} {
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
		return getPath(ret, path[1:]...)
	case map[string]interface{}:
		ret, ok := t[k]
		if !ok {
			return nil
		}
		return getPath(ret, path[1:]...)
	case []interface{}:
		i, err := strconv.Atoi(k)
		if err != nil {
			return nil
		}
		var ret interface{}
		if len(t) > i {
			ret = t[i]
		}
		return getPath(ret, path[1:]...)
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
				return errors.Errorf("unhandled [%T]", t)
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
				return errors.Errorf("unhandled [%T]", t)
			}
		}
	}
	return nil
}
