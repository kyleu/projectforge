package util

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/pkg/errors"
)

func (m ValueMap) ParseInt(path string) (int, error) {
	switch t := m.GetPath(path).(type) {
	case int:
		return t, nil
	case string:
		return strconv.Atoi(t)
	default:
		return 0, invalidTypeError(path, "int", t)
	}
}

func (m ValueMap) ParseMap(path string) (ValueMap, error) {
	switch t := m.GetPath(path).(type) {
	case ValueMap:
		return t, nil
	case map[string]interface{}:
		return t, nil
	case string:
		if strings.TrimSpace(t) == "" {
			return nil, nil
		}
		ret := ValueMap{}
		err := FromJSON([]byte(t), &ret)
		if err != nil {
			return nil, decorateError(m, path, "time", errors.Wrap(err, "invalid JSON"))
		}
		return ret, err
	default:
		return nil, invalidTypeError(path, "map", t)
	}
}

func (m ValueMap) ParseString(path string) (string, error) {
	return fmt.Sprint(m.GetPath(path)), nil
}

func (m ValueMap) ParseTime(path string) (time.Time, error) {
	ret, err := m.ParseTimeOpt(path)
	if err != nil {
		return time.Time{}, err
	}
	return *ret, nil
}

func (m ValueMap) ParseTimeOpt(path string) (*time.Time, error) {
	switch t := m.GetPath(path).(type) {
	case time.Time:
		return &t, nil
	case *time.Time:
		return t, nil
	case string:
		ret, err := TimeFromJS(t)
		return ret, decorateError(m, path, "time", err)
	default:
		return nil, invalidTypeError(path, "time", t)
	}
}

func (m ValueMap) ParseUUID(path string) (uuid.UUID, error) {
	ret, err := m.ParseUUIDOpt(path)
	if err != nil {
		return uuid.UUID{}, err
	}
	return *ret, nil
}

func (m ValueMap) ParseUUIDOpt(path string) (*uuid.UUID, error) {
	switch t := m.GetPath(path).(type) {
	case *uuid.UUID:
		return t, nil
	case uuid.UUID:
		return &t, nil
	case string:
		return UUIDFromString(t), nil
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

func invalidTypeError(path string, t string, observed interface{}) error {
	return errors.Errorf("unable to parse [%s] at path [%s], invalid type [%T]", t, path, observed)
}
