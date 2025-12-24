package util

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/samber/lo"
)

func (m ValueMap) ParseBool(path string, allowMissing bool, allowEmpty bool) (bool, error) {
	return parseMapField(m, path, allowMissing, func(res any) (bool, error) {
		return ParseBool(res, path, allowEmpty)
	})
}

func (m ValueMap) ParseFloat(path string, allowMissing bool, allowEmpty bool) (float64, error) {
	return parseMapField(m, path, allowMissing, func(res any) (float64, error) {
		return ParseFloat(res, path, allowEmpty)
	})
}

func (m ValueMap) ParseInt(path string, allowMissing bool, allowEmpty bool) (int, error) {
	return parseMapField(m, path, allowMissing, func(res any) (int, error) {
		return ParseInt(res, path, allowEmpty)
	})
}

func (m ValueMap) ParseInt64(path string, allowMissing bool, allowEmpty bool) (int64, error) {
	return parseMapField(m, path, allowMissing, func(res any) (int64, error) {
		return ParseInt64(res, path, allowEmpty)
	})
}

func (m ValueMap) ParseInt32(path string, allowMissing bool, allowEmpty bool) (int32, error) {
	return parseMapField(m, path, allowMissing, func(res any) (int32, error) {
		return ParseInt32(res, path, allowEmpty)
	})
}

func (m ValueMap) ParseInt16(path string, allowMissing bool, allowEmpty bool) (int16, error) {
	return parseMapField(m, path, allowMissing, func(res any) (int16, error) {
		return ParseInt16(res, path, allowEmpty)
	})
}

func (m ValueMap) ParseJSON(path string, allowMissing bool, allowEmpty bool) (any, error) {
	return parseMapField(m, path, allowMissing, func(res any) (any, error) {
		return ParseJSON(res, path, allowEmpty)
	})
}

func (m ValueMap) ParseMap(path string, allowMissing bool, allowEmpty bool) (ValueMap, error) {
	return parseMapField(m, path, allowMissing, func(res any) (ValueMap, error) {
		return ParseMap(res, path, allowEmpty)
	})
}

func (m ValueMap) ParseString(path string, allowMissing bool, allowEmpty bool) (string, error) {
	return parseMapField(m, path, allowMissing, func(res any) (string, error) {
		return ParseString(res, path, allowEmpty)
	})
}

func (m ValueMap) ParseTime(path string, allowMissing bool, allowEmpty bool) (*time.Time, error) {
	return parseMapField(m, path, allowMissing, func(res any) (*time.Time, error) {
		return ParseTime(res, path, allowEmpty)
	})
}

func (m ValueMap) ParseUUID(path string, allowMissing bool, allowEmpty bool) (*uuid.UUID, error) {
	return parseMapField(m, path, allowMissing, func(res any) (*uuid.UUID, error) {
		return ParseUUID(res, path, allowEmpty)
	})
}

func (m ValueMap) ParseArray(path string, allowMissing bool, allowEmpty bool, coerce bool) ([]any, error) {
	return parseMapField(m, path, allowMissing, func(res any) ([]any, error) {
		return ParseArray(res, path, allowEmpty, coerce)
	})
}

func (m ValueMap) ParseArrayString(path string, allowMissing bool, allowEmpty bool) ([]string, error) {
	return parseMapField(m, path, allowMissing, func(res any) ([]string, error) {
		return ParseArrayString(res, path, allowEmpty)
	})
}

func (m ValueMap) ParseArrayInt(path string, allowMissing bool, allowEmpty bool) ([]int, error) {
	return parseMapField(m, path, allowMissing, func(res any) ([]int, error) {
		return ParseArrayInt(res, path, allowEmpty)
	})
}

func (m ValueMap) ParseArrayFloat(path string, allowMissing bool, allowEmpty bool) ([]float64, error) {
	return parseMapField(m, path, allowMissing, func(res any) ([]float64, error) {
		return ParseArrayFloat(res, path, allowEmpty)
	})
}

func (m ValueMap) ParseArrayMap(path string, allowMissing bool, allowEmpty bool) ([]ValueMap, error) {
	return parseMapField(m, path, allowMissing, func(res any) ([]ValueMap, error) {
		return ParseArrayMap(res, path, allowEmpty)
	})
}

func parseMapField[T any](m ValueMap, path string, allowMissing bool, fn func(res any) (T, error)) (T, error) {
	result, err := m.GetPath(path, allowMissing)
	if err != nil {
		var x T
		return x, errors.Wrapf(err, "invalid %T", result)
	}
	return fn(result)
}

func simplifyValue(k string, v any) (any, error) {
	var err error
	if v != nil {
		switch t := v.(type) {
		case []string:
			v, err = ArrayFromAny[any](t)
		case []int:
			v, err = ArrayFromAny[any](t)
		case []int64:
			v, err = ArrayFromAny[any](t)
		case []float32:
			v, err = ArrayFromAny[any](t)
		case []float64:
			v, err = ArrayFromAny[any](t)
		case uuid.UUID:
			v = t.String()
		case *uuid.UUID:
			v = t.String()
		case time.Time:
			v = TimeToJS(&t)
		case *time.Time:
			v = TimeToJS(t)
		case ValueMap:
			v = t.AsMap(true)
		case []ValueMap:
			v = lo.Map(t, func(x ValueMap, _ int) any {
				return x.AsMap(true)
			})
		case []any:
			var acc []any
			for idx, x := range t {
				var simple any
				simple, err = simplifyValue(fmt.Sprintf("%s[%d]", k, idx), x)
				if err != nil {
					break
				}
				acc = append(acc, simple)
			}
			v = acc
		case map[string]any:
		case string:
		case int:
		case bool:
		case float64:
		default:
			v, err = FromJSONAny(ToJSONBytes(v, true))
			if err != nil {
				panic(fmt.Sprintf("encountered [%s] value of type [%T]", k, v))
			}
		}
	}
	return v, err
}
