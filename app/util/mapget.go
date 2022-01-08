package util

import (
	"encoding/json"
	"strings"
	"time"

	"github.com/pkg/errors"
)

func (m ValueMap) GetRequired(k string) (interface{}, error) {
	v, ok := m[k]
	if !ok {
		msg := "no value [%s] among candidates [%s]"
		return nil, errors.Errorf(msg, k, StringArrayOxfordComma(m.Keys(), "and"))
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

func (m ValueMap) GetInt64(k string, allowEmpty bool) (int64, error) {
	v, err := m.GetRequired(k)
	if err != nil {
		return 0, err
	}

	switch t := v.(type) {
	case int:
		return int64(t), nil
	case int32:
		return int64(t), nil
	case int64:
		return t, nil
	case nil:
		if allowEmpty {
			return 0, nil
		}
		return 0, errors.New(k + " is nil, not integer")
	default:
		return 0, errors.Errorf("expected integer, encountered %T", t)
	}
}

func (m ValueMap) GetFloat64(k string, allowEmpty bool) (float64, error) {
	v, err := m.GetRequired(k)
	if err != nil {
		return 0, err
	}

	switch t := v.(type) {
	case int:
		return float64(t), nil
	case int32:
		return float64(t), nil
	case float64:
		return t, nil
	case nil:
		if allowEmpty {
			return 0, nil
		}
		return 0, errors.New(k + " is nil, not float")
	default:
		return 0, errors.Errorf("expected float, encountered %T", t)
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

func (m ValueMap) GetTime(k string) (*time.Time, error) {
	v, err := m.GetRequired(k)
	if err != nil {
		return nil, err
	}

	var ret *time.Time
	switch t := v.(type) {
	case time.Time:
		ret = &t
	default:
		return nil, errors.Errorf("expected time, encountered %T", t)
	}
	return ret, nil
}

func (m ValueMap) GetType(k string, ret interface{}) error {
	v, err := m.GetRequired(k)
	if err != nil {
		return err
	}

	switch t := v.(type) {
	case []byte:
		err = json.Unmarshal(t, ret)
		if err != nil {
			return errors.Wrap(err, "failed to unmarshal to expected type")
		}
		return nil

	default:
		return errors.Errorf("expected binary json data, encountered %T", t)
	}
}

func (m ValueMap) GetMap(key string) (ValueMap, error) {
	switch t := m.GetPath(key).(type) {
	case map[string]interface{}:
		return t, nil
	case ValueMap:
		return t, nil
	case nil:
		return nil, nil
	default:
		return nil, errors.Errorf("unhandled type [%T] for key [%s], expected map", t, key)
	}
}

func (m ValueMap) GetArray(key string, allowEmpty bool) ([]interface{}, error) {
	switch t := m.GetPath(key).(type) {
	case []interface{}:
		return t, nil
	case nil:
		if !allowEmpty {
			return nil, errors.Errorf("could not find array for key [%s]", key)
		}
		return nil, nil
	default:
		return nil, errors.Errorf("unhandled type [%T] for key [%s], expected array", t, key)
	}
}
