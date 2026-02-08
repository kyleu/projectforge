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

func ParseBoolSimple(r any) bool {
	ret, _ := ParseBool(r, "", true)
	return ret
}

func ParseFloatSimple(r any) float64 {
	ret, _ := ParseFloat(r, "", true)
	return ret
}

func ParseIntSimple(r any) int {
	ret, _ := ParseInt(r, "", true)
	return ret
}

func ParseInt16Simple(r any) int16 {
	ret, _ := ParseInt16(r, "", true)
	return ret
}

func ParseInt32Simple(r any) int32 {
	ret, _ := ParseInt32(r, "", true)
	return ret
}

func ParseInt64Simple(r any) int64 {
	ret, _ := ParseInt64(r, "", true)
	return ret
}

func ParseJSONSimple(r any) any {
	ret, _ := ParseJSON(r, "", true)
	return ret
}

func ParseMapSimple(r any) ValueMap {
	ret, _ := ParseMap(r, "", true)
	return ret
}

func ParseOrderedMapSimple(r any) *OrderedMap[any] {
	ret, _ := ParseOrderedMap(r, "", true)
	return ret
}

func ParseStringSimple(r any) string {
	ret, _ := ParseString(r, "", true)
	return ret
}

func ParseTimeSimple(r any) *time.Time {
	ret, _ := ParseTime(r, "", true)
	return ret
}

func ParseUUIDSimple(r any) *uuid.UUID {
	ret, _ := ParseUUID(r, "", true)
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
	ret, err := ParseInt64(r, path, allowEmpty)
	if err != nil {
		return 0, err
	}
	return int(ret), nil
}

func ParseInt16(r any, path string, allowEmpty bool) (int16, error) {
	switch t := r.(type) {
	case int:
		return int16(t), nil
	case int16:
		return t, nil
	case int32:
		return int16(t), nil
	case int64:
		return int16(t), nil
	case float64:
		return int16(t), nil
	case string:
		ret, err := strconv.ParseInt(t, 10, 16)
		return int16(ret), err
	case []byte:
		ret, err := strconv.ParseInt(string(t), 10, 16)
		return int16(ret), err
	case nil:
		if !allowEmpty {
			return 0, errors.Errorf("could not find int16 for path [%s]", path)
		}
		return 0, nil
	default:
		return 0, invalidTypeError(path, "int16", t)
	}
}

func ParseInt32(r any, path string, allowEmpty bool) (int32, error) {
	switch t := r.(type) {
	case int:
		return int32(t), nil
	case int32:
		return t, nil
	case int64:
		return int32(t), nil
	case float64:
		return int32(t), nil
	case string:
		ret, err := strconv.ParseInt(t, 10, 32)
		return int32(ret), err
	case []byte:
		ret, err := strconv.ParseInt(string(t), 10, 32)
		return int32(ret), err
	case nil:
		if !allowEmpty {
			return 0, errors.Errorf("could not find int32 for path [%s]", path)
		}
		return 0, nil
	default:
		return 0, invalidTypeError(path, "int32", t)
	}
}

func ParseInt64(r any, path string, allowEmpty bool) (int64, error) {
	switch t := r.(type) {
	case int:
		return int64(t), nil
	case int32:
		return int64(t), nil
	case int64:
		return t, nil
	case float64:
		return int64(t), nil
	case string:
		ret, err := strconv.ParseInt(t, 10, 64)
		return ret, err
	case []byte:
		ret, err := strconv.ParseInt(string(t), 10, 64)
		return ret, err
	case nil:
		if !allowEmpty {
			return 0, errors.Errorf("could not find int for path [%s]", path)
		}
		return 0, nil
	default:
		return 0, invalidTypeError(path, "int", t)
	}
}

func ParseJSON(r any, path string, allowEmpty bool) (any, error) {
	switch t := r.(type) {
	case []byte:
		return FromJSONAny(t)
	case string:
		return FromJSONAny([]byte(t))
	case nil:
		if !allowEmpty {
			return nil, errors.Errorf("could not find json for path [%s]", path)
		}
		return nil, nil
	default:
		return t, nil
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
		return nil, invalidTypeError(path, KeyMap, t)
	}
}

func ParseOrderedMap(r any, path string, allowEmpty bool) (*OrderedMap[any], error) {
	switch t := r.(type) {
	case *OrderedMap[any]:
		if (!allowEmpty) && len(t.Map) == 0 {
			return nil, errors.New("empty map")
		}
		return t, nil
	case ValueMap:
		if (!allowEmpty) && len(t) == 0 {
			return nil, errors.New("empty map")
		}
		o := NewOrderedMap[any](false, len(t))
		for k, v := range t {
			o.Set(k, v)
		}
		return o, nil
	case map[string]any:
		if (!allowEmpty) && len(t) == 0 {
			return nil, errors.New("empty map")
		}
		o := NewOrderedMap[any](false, len(t))
		for k, v := range t {
			o.Set(k, v)
		}
		return o, nil
	case string:
		if strings.TrimSpace(t) == "" {
			return nil, nil
		}
		ret, err := FromJSONOrderedMap[any]([]byte(t))
		if err != nil {
			return nil, wrapError(path, "time", errors.Wrap(err, "invalid JSON"))
		}
		return ret, err
	case []byte:
		if len(t) == 0 {
			return nil, nil
		}
		ret, err := FromJSONOrderedMap[any](t)
		if err != nil {
			return nil, wrapError(path, "time", errors.Wrap(err, "invalid JSON"))
		}
		return ret, err
	case nil:
		if !allowEmpty {
			return nil, errors.Errorf("could not find ordered map for path [%s]", path)
		}
		return nil, nil
	default:
		return nil, invalidTypeError(path, "ordered map", t)
	}
}

func ParseString(r any, path string, allowEmpty bool) (string, error) {
	switch t := r.(type) {
	case rune:
		return string(t), nil
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
		return StringJoin(t, ","), nil
	case map[string]any:
		return ToJSONCompact(t), nil
	case ValueMap:
		return ToJSONCompact(t), nil
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
