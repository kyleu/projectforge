package util

import (
	"time"

	"github.com/google/uuid"
	"github.com/pkg/errors"
)

func (m ValueMap) GetRequired(k string) (any, error) {
	v, ok := m[k]
	if !ok {
		const msg = "no value [%s] among candidates [%s]"
		return nil, errors.Errorf(msg, k, StringArrayOxfordComma(m.Keys(), "and"))
	}
	return v, nil
}

func (m ValueMap) GetArray(key string, allowEmpty bool) ([]any, error) {
	return m.ParseArray(key, false, allowEmpty, false)
}

func (m ValueMap) GetArrayOpt(key string) []any {
	ret, _ := m.ParseArray(key, false, true, false)
	return ret
}

func (m ValueMap) GetArrayOr(key string, dflt []any, allowEmpty bool) []any {
	if ret, err := m.GetArray(key, allowEmpty); err == nil {
		return ret
	}
	return dflt
}

func (m ValueMap) GetStringArray(key string, allowEmpty bool) ([]string, error) {
	return m.ParseArrayString(key, false, allowEmpty)
}

func (m ValueMap) GetStringArrayOpt(key string) []string {
	ret, _ := m.ParseArrayString(key, false, true)
	return ret
}

func (m ValueMap) GetStringArrayOr(key string, dflt []string, allowEmpty bool) []string {
	if ret, err := m.GetStringArray(key, allowEmpty); err == nil {
		return ret
	}
	return dflt
}

func (m ValueMap) GetIntArray(key string, allowEmpty bool) ([]int, error) {
	return m.ParseArrayInt(key, false, allowEmpty)
}

func (m ValueMap) GetIntArrayOpt(key string) []int {
	ret, _ := m.ParseArrayInt(key, false, true)
	return ret
}

func (m ValueMap) GetIntArrayOr(key string, dflt []int, allowEmpty bool) []int {
	if ret, err := m.GetIntArray(key, allowEmpty); err == nil {
		return ret
	}
	return dflt
}

func (m ValueMap) GetFloatArray(key string, allowEmpty bool) ([]float64, error) {
	return m.ParseArrayFloat(key, false, allowEmpty)
}

func (m ValueMap) GetFloatArrayOpt(key string) []float64 {
	ret, _ := m.ParseArrayFloat(key, false, true)
	return ret
}

func (m ValueMap) GetFloatArrayOr(key string, dflt []float64, allowEmpty bool) []float64 {
	if ret, err := m.GetFloatArray(key, allowEmpty); err == nil {
		return ret
	}
	return dflt
}

func (m ValueMap) GetMapArray(key string, allowEmpty bool) ([]ValueMap, error) {
	return m.ParseArrayMap(key, false, allowEmpty)
}

func (m ValueMap) GetMapArrayOpt(key string) []ValueMap {
	ret, _ := m.ParseArrayMap(key, false, true)
	return ret
}

func (m ValueMap) GetMapArrayOr(key string, dflt []ValueMap, allowEmpty bool) []ValueMap {
	if ret, err := m.GetMapArray(key, allowEmpty); err == nil {
		return ret
	}
	return dflt
}

func (m ValueMap) GetBool(key string, allowEmpty bool) (bool, error) {
	return m.ParseBool(key, false, allowEmpty)
}

func (m ValueMap) GetBoolOpt(key string) bool {
	ret, _ := m.ParseBool(key, true, true)
	return ret
}

func (m ValueMap) GetBoolOr(key string, dflt bool, allowEmpty bool) bool {
	if ret, err := m.GetBool(key, allowEmpty); err == nil {
		return ret
	}
	return dflt
}

func (m ValueMap) GetInt(key string, allowEmpty bool) (int, error) {
	return m.ParseInt(key, false, allowEmpty)
}

func (m ValueMap) GetIntOpt(key string) int {
	ret, _ := m.ParseInt(key, false, true)
	return ret
}

func (m ValueMap) GetIntOr(key string, dflt int, allowEmpty bool) int {
	if ret, err := m.GetInt(key, allowEmpty); err == nil {
		return ret
	}
	return dflt
}

func (m ValueMap) GetInt64(key string, allowEmpty bool) (int64, error) {
	ret, err := m.ParseInt(key, false, allowEmpty)
	return int64(ret), err
}

func (m ValueMap) GetInt64Opt(key string) int64 {
	ret, _ := m.ParseInt(key, false, true)
	return int64(ret)
}

func (m ValueMap) GetInt64Or(key string, dflt int64, allowEmpty bool) int64 {
	if ret, err := m.GetInt64(key, allowEmpty); err == nil {
		return ret
	}
	return dflt
}

func (m ValueMap) GetFloat(key string, allowEmpty bool) (float64, error) {
	return m.ParseFloat(key, false, allowEmpty)
}

func (m ValueMap) GetFloatOpt(key string) float64 {
	ret, _ := m.ParseFloat(key, false, true)
	return ret
}

func (m ValueMap) GetFloatOr(key string, dflt float64, allowEmpty bool) float64 {
	if ret, err := m.GetFloat(key, allowEmpty); err == nil {
		return ret
	}
	return dflt
}

func (m ValueMap) GetMap(key string, allowEmpty bool) (ValueMap, error) {
	return m.ParseMap(key, false, allowEmpty)
}

func (m ValueMap) GetMapOpt(key string) ValueMap {
	ret, _ := m.ParseMap(key, false, true)
	return ret
}

func (m ValueMap) GetMapOr(key string, dflt ValueMap, allowEmpty bool) ValueMap {
	if ret, err := m.GetMap(key, allowEmpty); err == nil {
		return ret
	}
	return dflt
}

func (m ValueMap) GetString(key string, allowEmpty bool) (string, error) {
	return m.ParseString(key, false, allowEmpty)
}

func (m ValueMap) GetStringOpt(key string) string {
	ret, _ := m.ParseString(key, true, true)
	return ret
}

func (m ValueMap) GetStringOr(key string, dflt string, allowEmpty bool) string {
	ret, err := m.GetString(key, allowEmpty)
	if err != nil {
		return dflt
	}
	return ret
}

func (m ValueMap) GetStringPtr(key string) *string {
	ret, err := m.ParseString(key, true, true)
	if err != nil {
		return nil
	}
	return &ret
}

func (m ValueMap) GetRichString(key string, allowEmpty bool) (RichString, error) {
	ret, err := m.ParseString(key, false, allowEmpty)
	return RS(ret), err
}

func (m ValueMap) GetRichStringOpt(key string) RichString {
	ret, _ := m.GetRichString(key, true)
	return ret
}

func (m ValueMap) GetTime(key string, allowEmpty bool) (*time.Time, error) {
	return m.ParseTime(key, false, allowEmpty)
}

func (m ValueMap) GetTimeOpt(key string) time.Time {
	ret, _ := m.ParseTime(key, true, true)
	if ret == nil {
		return time.Time{}
	}
	return *ret
}

func (m ValueMap) GetTimeOr(key string, dflt *time.Time, allowEmpty bool) *time.Time {
	if ret, err := m.GetTime(key, allowEmpty); err == nil {
		return ret
	}
	return dflt
}

func (m ValueMap) GetUUID(key string, allowEmpty bool) (*uuid.UUID, error) {
	return m.ParseUUID(key, false, allowEmpty)
}

func (m ValueMap) GetUUIDOpt(key string) uuid.UUID {
	ret, _ := m.ParseUUID(key, true, true)
	if ret == nil {
		return uuid.UUID{}
	}
	return *ret
}

func (m ValueMap) GetUUIDOr(key string, dflt *uuid.UUID, allowEmpty bool) *uuid.UUID {
	if ret, err := m.GetUUID(key, allowEmpty); err == nil {
		return ret
	}
	return dflt
}

func (m ValueMap) GetType(key string, ret any) error {
	result, err := m.GetPath(key, false)
	if err != nil {
		return errors.Wrap(err, "invalid type")
	}
	switch t := result.(type) {
	case []byte:
		err = FromJSON(t, ret)
		if err != nil {
			return errors.Wrap(err, "failed to unmarshal to expected type")
		}
		return nil
	default:
		return errors.Errorf("expected binary json data, encountered %T", t)
	}
}

func MapGetOrElse[K comparable, V any](mp map[K]V, k K, dflt V) V {
	ret, ok := mp[k]
	if ok {
		return ret
	}
	return dflt
}
