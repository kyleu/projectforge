package util

import (
	"time"

	"github.com/google/uuid"
	"github.com/pkg/errors"
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
