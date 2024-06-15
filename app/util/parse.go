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

func DefaultValue[T any]() T {
	var ret T
	return ret
}

func ParseBool(r any, path string, allowEmpty bool) (bool, error) {
	switch t := r.(type) {
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

func ParseFloat(r any, path string, allowEmpty bool) (float64, error) {
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

func ParseInt(r any, path string, allowEmpty bool) (int, error) {
	switch t := r.(type) {
	case int:
		return t, nil
	case int64:
		return int(t), nil
	case float64:
		return int(t), nil
	case string:
		ret, err := strconv.ParseInt(t, 10, 64)
		return int(ret), err
	case []byte:
		ret, err := strconv.ParseInt(string(t), 10, 64)
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

func ParseMap(r any, path string, allowEmpty bool) (ValueMap, error) {
	switch t := r.(type) {
	case ValueMap:
		if (!allowEmpty) && len(t) == 0 {
			return nil, errors.New("empty map")
		}
		return t, nil
	case map[string]any:
		if (!allowEmpty) && len(t) == 0 {
			return nil, errors.New("empty map")
		}
		return t, nil
	case string:
		if strings.TrimSpace(t) == "" {
			return nil, nil
		}
		ret, err := FromJSONMap([]byte(t))
		if err != nil {
			return nil, wrapError(path, "time", errors.Wrap(err, "invalid JSON"))
		}
		return ret, err
	case []byte:
		if len(t) == 0 {
			return nil, nil
		}
		ret, err := FromJSONMap(t)
		if err != nil {
			return nil, wrapError(path, "time", errors.Wrap(err, "invalid JSON"))
		}
		return ret, err
	case nil:
		if !allowEmpty {
			return nil, errors.Errorf("could not find map for path [%s]", path)
		}
		return nil, nil
	default:
		return nil, invalidTypeError(path, "map", t)
	}
}

func ParseString(r any, path string, allowEmpty bool) (string, error) {
	switch t := r.(type) {
	case string:
		if (!allowEmpty) && t == "" {
			return "", errors.New("empty string")
		}
		return t, nil
	case []byte:
		if (!allowEmpty) && len(t) == 0 {
			return "", errors.New("empty string")
		}
		return string(t), nil
	case []string:
		if (!allowEmpty) && (len(t) == 0 || t[0] == "") {
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

func ParseTime(r any, path string, allowEmpty bool) (*time.Time, error) {
	switch t := r.(type) {
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
			return nil, wrapError(path, t, err)
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

func ParseUUID(r any, path string, allowEmpty bool) (*uuid.UUID, error) {
	switch t := r.(type) {
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
	case []byte:
		if len(t) == 16 {
			ret, err := uuid.FromBytes(t)
			if err != nil {
				return nil, err
			}
			if ret == uuid.Nil && (!allowEmpty) {
				return nil, errors.Errorf("could not parse uuid from path [%s]", path)
			}
			return &ret, nil
		}
		return nil, errors.Errorf("invalid uuid bytes with length [%d]", len(t))
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

func invalidTypeError(path string, t string, observed any) error {
	return errors.Errorf("unable to parse [%s] at path [%s], invalid type [%T]", t, path, observed)
}

func wrapError(path string, t string, err error) error {
	if err == nil {
		return nil
	}
	return errors.Wrapf(err, "error parsing [%s] as [%s]", path, t)
}
