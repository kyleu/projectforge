// Package util - Content managed by Project Forge, see [projectforge.md] for details.
package util

import (
	"strings"

	"github.com/pkg/errors"
)

func (m ValueMap) ParseMap(path string, allowMissing bool, allowEmpty bool) (ValueMap, error) {
	result, err := m.GetPath(path, allowMissing)
	if err != nil {
		return nil, errors.Wrap(err, "invalid int")
	}
	switch t := result.(type) {
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
			return nil, decorateError(m, path, "time", errors.Wrap(err, "invalid JSON"))
		}
		return ret, err
	case []byte:
		if len(t) == 0 {
			return nil, nil
		}
		ret, err := FromJSONMap(t)
		if err != nil {
			return nil, decorateError(m, path, "time", errors.Wrap(err, "invalid JSON"))
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
