// Package util - Content managed by Project Forge, see [projectforge.md] for details.
package util

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/pkg/errors"
)

func (m ValueMap) ParseBool(path string, allowMissing bool, allowEmpty bool) (bool, error) {
	result, err := m.GetPath(path, allowMissing)
	if err != nil {
		return false, errors.Wrap(err, "invalid bool")
	}
	switch t := result.(type) {
	case bool:
		return t, nil
	case string:
		return t == BoolTrue, nil
	case nil:
		if !allowEmpty {
			return false, errors.Errorf("could not find bool for path [%s]", path)
		}
		return false, nil
	default:
		return false, invalidTypeError(path, "bool", t)
	}
}

func (m ValueMap) ParseFloat(path string, allowMissing bool, allowEmpty bool) (float64, error) {
	result, err := m.GetPath(path, allowMissing)
	if err != nil {
		return 0, errors.Wrap(err, "invalid float")
	}
	switch t := result.(type) {
	case int:
		return float64(t), nil
	case int64:
		return float64(t), nil
	case float64:
		return t, nil
	case string:
		return strconv.ParseFloat(t, 64)
	case nil:
		if !allowEmpty {
			return 0, errors.Errorf("could not find float for path [%s]", path)
		}
		return 0, nil
	default:
		return 0, invalidTypeError(path, "float", t)
	}
}

func (m ValueMap) ParseInt(path string, allowMissing bool, allowEmpty bool) (int, error) {
	result, err := m.GetPath(path, allowMissing)
	if err != nil {
		return 0, errors.Wrap(err, "invalid int")
	}
	return valueInt(path, result, allowEmpty)
}

func (m ValueMap) ParseString(path string, allowMissing bool, allowEmpty bool) (string, error) {
	result, err := m.GetPath(path, allowMissing)
	if err != nil {
		return "", errors.Wrap(err, "invalid string")
	}
	switch t := result.(type) {
	case string:
		if (!allowEmpty) && t == "" {
			return "", errors.New("empty string")
		}
		return t, nil
	case []string:
		if (!allowEmpty) && len(t) == 0 || t[0] == "" {
			return "", errors.New("empty string")
		}
		return strings.Join(t, "||"), nil
	case nil:
		if !allowEmpty {
			return "", errors.Errorf("could not find string for path [%s]", path)
		}
		return "", nil
	default:
		return fmt.Sprint(t), nil
	}
}

func (m ValueMap) ParseTime(path string, allowMissing bool, allowEmpty bool) (*time.Time, error) {
	result, err := m.GetPath(path, allowMissing)
	if err != nil {
		return nil, errors.Wrap(err, "invalid time")
	}
	switch t := result.(type) {
	case time.Time:
		return &t, nil
	case *time.Time:
		if t == nil && (!allowEmpty) {
			return nil, errors.New("empty time")
		}
		return t, nil
	case string:
		ret, err := TimeFromString(t)
		if err != nil {
			return nil, decorateError(m, path, "time", err)
		}
		if ret == nil && (!allowEmpty) {
			return nil, errors.Errorf("invalid time [%s]", t)
		}
		return ret, nil
	case nil:
		if !allowEmpty {
			return nil, errors.Errorf("could not find time for path [%s]", path)
		}
		return nil, nil
	default:
		return nil, invalidTypeError(path, "time", t)
	}
}

func (m ValueMap) ParseUUID(path string, allowMissing bool, allowEmpty bool) (*uuid.UUID, error) {
	result, err := m.GetPath(path, allowMissing)
	if err != nil {
		return nil, errors.Wrap(err, "invalid uuid")
	}

	switch t := result.(type) {
	case *uuid.UUID:
		if t == nil && (!allowEmpty) {
			return nil, errors.New("empty uuid")
		}
		return t, nil
	case uuid.UUID:
		if t == uuid.Nil && (!allowEmpty) {
			return nil, errors.New("empty uuid")
		}
		return &t, nil
	case string:
		if t == "" && allowEmpty {
			return nil, nil
		}
		ret, err := uuid.Parse(t)
		if err != nil {
			return nil, err
		}
		if ret == uuid.Nil && (!allowEmpty) {
			return nil, errors.Errorf("could not find uuid for path [%s]", path)
		}
		return &ret, nil
	case nil:
		if !allowEmpty {
			return nil, errors.Errorf("could not find uuid for path [%s]", path)
		}
		return nil, nil
	default:
		return nil, invalidTypeError(path, "uuid", t)
	}
}

func decorateError(m ValueMap, path string, t string, err error) error {
	if err == nil {
		return nil
	}
	return errors.Wrapf(err, "error parsing [%s] as [%s] from map with fields [%s]", path, t, strings.Join(m.Keys(), ", "))
}

func invalidTypeError(path string, t string, observed any) error {
	return errors.Errorf("unable to parse [%s] at path [%s], invalid type [%T]", t, path, observed)
}

func valueInt(path string, r any, allowEmpty bool) (int, error) {
	switch t := r.(type) {
	case int:
		return t, nil
	case int64:
		return int(t), nil
	case float64:
		return int(t), nil
	case string:
		ret, err := strconv.ParseInt(t, 10, 32)
		return int(ret), err
	case nil:
		if !allowEmpty {
			return 0, errors.Errorf("could not find int for path [%s]", path)
		}
		return 0, nil
	default:
		return 0, invalidTypeError(path, "int", t)
	}
}

func valueFloat(path string, r any, allowEmpty bool) (float64, error) {
	switch t := r.(type) {
	case int:
		return float64(t), nil
	case int64:
		return float64(t), nil
	case float64:
		return t, nil
	case string:
		ret, err := strconv.ParseFloat(t, 32)
		return ret, err
	case nil:
		if !allowEmpty {
			return 0, errors.Errorf("could not find float for path [%s]", path)
		}
		return 0, nil
	default:
		return 0, invalidTypeError(path, "float", t)
	}
}

func valueMap(path string, r any, allowEmpty bool) (ValueMap, error) {
	switch t := r.(type) {
	case ValueMap:
		return t, nil
	case map[string]any:
		return t, nil
	case string:
		return FromJSONMap([]byte(t))
	case []byte:
		return FromJSONMap(t)
	case nil:
		if !allowEmpty {
			return nil, errors.Errorf("could not find int for path [%s]", path)
		}
		return ValueMap{}, nil
	default:
		return nil, invalidTypeError(path, "int", t)
	}
}
