package util

import (
	"fmt"
	"strings"

	"github.com/pkg/errors"
	"github.com/samber/lo"
)

func ParseArray(r any, path string, allowEmpty bool, coerce bool) ([]any, error) {
	switch t := r.(type) {
	case string:
		if strings.TrimSpace(t) == "" {
			return nil, nil
		}
		var ret []any
		err := FromJSON([]byte(t), &ret)
		if err != nil {
			if coerce {
				return lo.ToAnySlice(StringSplitAndTrim(t, ",")), nil
			}
			return nil, wrapError(path, KeyJSON, errors.Wrap(err, "invalid JSON string"))
		}
		return ret, err
	case []byte:
		if len(t) == 0 {
			return nil, nil
		}
		var ret []any
		err := FromJSON(t, &ret)
		if err != nil {
			if coerce {
				return lo.ToAnySlice(StringSplitAndTrim(string(t), ",")), nil
			}
			return nil, wrapError(path, KeyJSON, errors.Wrap(err, "invalid JSON bytes"))
		}
		return ret, err
	case []any:
		if (!allowEmpty) && len(t) == 0 {
			return nil, errors.New("empty array")
		}
		return t, nil
	case []string:
		if (!allowEmpty) && len(t) == 0 {
			return nil, errors.New("empty array")
		}
		return lo.ToAnySlice(t), nil
	case []int:
		if (!allowEmpty) && len(t) == 0 {
			return nil, errors.New("empty array")
		}
		return lo.ToAnySlice(t), nil
	case nil:
		if !allowEmpty {
			return nil, errors.Errorf("could not find array for path [%s]", path)
		}
		return nil, nil
	default:
		return nil, invalidTypeError(path, "array", t)
	}
}

func ParseArrayString(r any, path string, allowEmpty bool) ([]string, error) {
	a, err := ParseArray(r, path, allowEmpty, true)
	if err != nil {
		return nil, err
	}
	return lo.Map(a, func(x any, _ int) string {
		return fmt.Sprint(x)
	}), nil
}

func ParseArrayInt(r any, path string, allowEmpty bool) ([]int, error) {
	a, err := ParseArray(r, path, allowEmpty, true)
	if err != nil {
		return nil, err
	}
	ia := make([]int, 0, len(a))
	for idx, x := range a {
		i, err := ParseInt(x, fmt.Sprintf("%s.%d", path, idx), allowEmpty)
		if err != nil {
			return nil, err
		}
		ia = append(ia, i)
	}
	return ia, nil
}

func ParseArrayFloat(r any, path string, allowEmpty bool) ([]float64, error) {
	a, err := ParseArray(r, path, allowEmpty, true)
	if err != nil {
		return nil, err
	}
	fa := make([]float64, 0, len(a))
	for idx, x := range a {
		f, err := ParseFloat(x, fmt.Sprintf("%s.%d", path, idx), allowEmpty)
		if err != nil {
			return nil, err
		}
		fa = append(fa, f)
	}
	return fa, nil
}

func ParseArrayMap(r any, path string, allowEmpty bool) ([]ValueMap, error) {
	a, err := ParseArray(r, path, allowEmpty, false)
	if err != nil {
		return nil, err
	}
	ma := make([]ValueMap, 0, len(a))
	for idx, x := range a {
		m, err := ParseMap(x, fmt.Sprintf("%s.%d", path, idx), allowEmpty)
		if err != nil {
			return nil, err
		}
		ma = append(ma, m)
	}
	return ma, nil
}
