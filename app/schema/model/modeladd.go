package model

import (
	"fmt"

	"github.com/kyleu/projectforge/app/schema/field"
	"github.com/pkg/errors"
)

func (m *Model) AddField(f *field.Field) error {
	if f == nil {
		return errors.New("nil field")
	}
	if _, v := m.Fields.Get(f.Key); v != nil {
		return errors.Errorf("field [%s] already exists", f.Key)
	}
	m.Fields = append(m.Fields, f)
	return nil
}

func (m *Model) AddIndex(i *Index) error {
	if i == nil {
		return errors.New("nil index")
	}
	if m.Indexes.Get(i.Key) != nil {
		return errors.Errorf("index [%s] already exists", i.Key)
	}
	m.Indexes = append(m.Indexes, i)
	return nil
}

func (m *Model) AddRelationship(r *Relationship) error {
	if r == nil {
		return errors.New("nil relation")
	}
	if m.Relationships.Get(r.Key) != nil {
		return errors.Errorf("relation [%s] already exists", r.Key)
	}
	m.Relationships = append(m.Relationships, r)
	return nil
}

func (m *Model) AddReference(r *Reference) error {
	if r == nil {
		return errors.New("nil reference")
	}
	if m.References.Get(r.Key) != nil {
		idx := 2
		test := fmt.Sprintf("%s-%d", r.Key, idx)
		for m.References.Get(test) != nil {
			idx++
			test = fmt.Sprintf("%s-%d", r.Key, idx)
		}
		r.Key = test
		return m.AddReference(r)
	}
	m.References = append(m.References, r)
	return nil
}
