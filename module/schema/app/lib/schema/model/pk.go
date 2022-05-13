package model

import (
	"github.com/pkg/errors"

	"{{{ .Package }}}/app/lib/schema/field"
	"{{{ .Package }}}/app/util"
)

func (m *Model) GetPK(logger util.Logger) []string {
	if m.pk == nil {
		for _, idx := range m.Indexes {
			if idx.Primary {
				if m.pk != nil && logger != nil {
					logger.Error("multiple primary keys?!")
				}
				m.pk = idx.Fields
			}
		}
	}
	return m.pk
}

func (m *Model) IsPK(key string, logger util.Logger) bool {
	pk := m.GetPK(logger)
	for _, col := range pk {
		if col == key {
			return true
		}
	}
	return false
}

func GetValues(src field.Fields, tgt []string, vals []any) ([]any, error) {
	if len(src) != len(vals) {
		return nil, errors.Errorf("[%d] fields provided, but [%d] values expected", len(vals), len(src))
	}
	ret := make([]any, 0, len(tgt))
	for _, t := range tgt {
		for idx, f := range src {
			if f.Key == t {
				ret = append(ret, vals[idx])
				break
			}
		}
	}
	return ret, nil
}

func GetStrings(src field.Fields, tgt []string, vals []any) ([]string, error) {
	is, err := GetValues(src, tgt, vals)
	if err != nil {
		return nil, err
	}
	ret := util.StringArrayFromInterfaces(is, 0)
	return ret, nil
}
