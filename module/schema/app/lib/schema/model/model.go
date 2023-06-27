package model

import (
	"fmt"

	"github.com/pkg/errors"
	"github.com/samber/lo"

	"{{{ .Package }}}/app/lib/schema/field"
	"{{{ .Package }}}/app/util"
)

type Model struct {
	Key           string        `json:"key"`
	Pkg           util.Pkg      `json:"pkg,omitempty"`
	Type          Type          `json:"type"`
	Title         string        `json:"-"` // override only
	Plural        string        `json:"-"` // override only
	Interfaces    []string      `json:"interfaces,omitempty"`
	Fields        field.Fields  `json:"fields,omitempty"`
	Indexes       Indexes       `json:"indexes,omitempty"`
	Relationships Relationships `json:"relationships,omitempty"`
	References    References    `json:"-"` // internal cache
	Metadata      *Metadata     `json:"metadata,omitempty"`
	pk            []string      // internal cache
}

func NewModel(pkg util.Pkg, key string) *Model {
	return &Model{Key: key, Pkg: pkg}
}

func (m *Model) String() string {
	if len(m.Pkg) == 0 {
		return m.Key
	}
	return m.Pkg.ToPath(m.Key)
}

func (m *Model) Name() string {
	if m.Title == "" {
		return util.StringToSingular(util.StringToTitle(m.Key))
	}
	return m.Title
}

func (m *Model) PluralName() string {
	if m.Plural == "" {
		ret := m.Name()
		return util.StringToPlural(ret)
	}
	return m.Plural
}

func (m *Model) Description() string {
	return fmt.Sprintf("%s model [%s]", m.Type.String(), m.Key)
}

func (m *Model) Path() util.Pkg {
	return m.Pkg.With(m.Key)
}

func (m *Model) PathString() string {
	return m.Pkg.ToPath(m.Key)
}

func (m *Model) OrderedMap(data []any) (*util.OrderedMap[any], error) {
	if len(data) != len(m.Fields) {
		return nil, errors.Errorf("expected [%d] data elements, found [%d]", len(m.Fields), len(data))
	}
	ret := util.NewOrderedMap[any](false, len(m.Fields))
	lo.ForEach(m.Fields, func(f *field.Field, idx int) {
		ret.Append(f.Key, data[idx])
	})
	return ret, nil
}
