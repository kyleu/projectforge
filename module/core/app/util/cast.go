package util

import "github.com/pkg/errors"

func Cast[T any](v any) (T, error) {
	ret, ok := v.(T)
	if !ok {
		return ret, errors.Errorf("unable to cast [%T] to [%T]", v, ret)
	}
	return ret, nil
}

func CastOK[T any](v any, loggers ...Logger) T {
	ret, err := Cast[T](v)
	if err != nil {
		for _, logger := range loggers {
			logger.Errorf("unable to cast [%T] to [%T]: %+v", v, ret, err)
		}
	}
	return ret
}
