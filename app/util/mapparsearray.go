// Package util - Content managed by Project Forge, see [projectforge.md] for details.
package util

import (
	"fmt"
	"strings"

	"github.com/pkg/errors"
	"github.com/samber/lo"
)

func (m ValueMap) ParseArray(path string, allowMissing bool, allowEmpty bool) ([]any, error) {
	result, err := m.GetPath(path, allowMissing)
	if err != nil {
		return nil, errors.Wrap(err, "invalid array")
	}
	switch t := result.(type) {
	case string:
		if strings.TrimSpace(t) == "" {
			return nil, nil
		}
		var ret []any
		err := FromJSON([]byte(t), &ret)
		if err != nil {
			return nil, decorateError(m, path, "time", errors.Wrap(err, "invalid JSON"))
		}
		return ret, err
	case []byte:
		if len(t) == 0 {
			return nil, nil
		}
		var ret []any
		err := FromJSON(t, &ret)
		if err != nil {
			return nil, decorateError(m, path, "time", errors.Wrap(err, "invalid JSON"))
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
		return InterfaceArrayFrom(t...), nil
	case []int:
		if (!allowEmpty) && len(t) == 0 {
			return nil, errors.New("empty array")
		}
		return InterfaceArrayFrom(t...), nil
	case nil:
		if !allowEmpty {
			return nil, errors.Errorf("could not find array for path [%s]", path)
		}
		return nil, nil
	default:
		return nil, invalidTypeError(path, "array", t)
	}
}

func (m ValueMap) ParseArrayInt(path string, allowMissing bool, allowEmpty bool) ([]int, error) {
	a, err := m.ParseArray(path, allowMissing, allowEmpty)
	if err != nil {
		return nil, err
	}
	ia := make([]int, 0, len(a))
	for idx, x := range a {
		i, err := valueInt(fmt.Sprintf("%s.%d", path, idx), x, allowEmpty)
		if err != nil {
			return nil, err
		}
		ia = append(ia, i)
	}
	return ia, nil
}

func (m ValueMap) ParseArrayString(path string, allowMissing bool, allowEmpty bool) ([]string, error) {
	a, err := m.ParseArray(path, allowMissing, allowEmpty)
	if err != nil {
		return nil, err
	}
	return lo.Map(a, func(x any, _ int) string {
		return fmt.Sprint(x)
	}), nil
}
