package model

import (
	"github.com/kyleu/projectforge/app/schema/field"
	"github.com/pkg/errors"
	"go.uber.org/zap"

	"github.com/kyleu/projectforge/app/util"
)

func (m *Model) GetPK(logger *zap.SugaredLogger) []string {
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

func (m *Model) IsPK(key string, logger *zap.SugaredLogger) bool {
	pk := m.GetPK(logger)
	for _, col := range pk {
		if col == key {
			return true
		}
	}
	return false
}

func GetValues(src field.Fields, tgt []string, vals []interface{}) ([]interface{}, error) {
	if len(src) != len(vals) {
		return nil, errors.Errorf("[%d] fields provided, but [%d] values expected", len(vals), len(src))
	}
	ret := make([]interface{}, 0, len(tgt))
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

func GetStrings(src field.Fields, tgt []string, vals []interface{}) ([]string, error) {
	is, err := GetValues(src, tgt, vals)
	if err != nil {
		return nil, err
	}
	ret := util.StringArrayFromInterfaces(is, 0)
	return ret, nil
}
