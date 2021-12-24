package schema

import (
	"fmt"
	"sort"
	"strings"

	"{{{ .Package }}}/app/schema/field"
	"{{{ .Package }}}/app/schema/model"
	"{{{ .Package }}}/app/util"
	"github.com/pkg/errors"
)

const (
	KeyTitle  = "title"
	KeyPlural = "plural"
)

type Override struct {
	Type string      `json:"type,omitempty"`
	Path util.Pkg    `json:"path,omitempty"`
	Prop string      `json:"prop,omitempty"`
	Val  interface{} `json:"val,omitempty"`
}

func NewOverride(typ string, path util.Pkg, prop string, val interface{}) *Override {
	return &Override{Type: typ, Path: path, Prop: prop, Val: val}
}

func (o Override) String() string {
	return fmt.Sprintf("%s[%s].%s = %v", o.Type, strings.Join(o.Path, "/"), o.Prop, o.Val)
}

type Overrides []*Override

func (o Overrides) Purge(path util.Pkg) Overrides {
	ret := make(Overrides, 0, len(o))
	for _, x := range o {
		if !(x.Path.StartsWith(path)) {
			ret = append(ret, x)
		}
	}
	return ret
}

func (o Overrides) Sort() {
	sort.Slice(o, func(i, j int) bool {
		l := o[i]
		r := o[j]
		if !l.Path.Equals(r.Path) {
			for idx, p := range l.Path {
				if idx >= len(r.Path) {
					return false
				}
				if p != r.Path[idx] {
					return p < r.Path[idx]
				}
			}
		}
		if len(l.Path) != len(r.Path) {
			return false
		}
		if l.Prop != r.Prop {
			return l.Prop < r.Prop
		}
		if l.Type != r.Type {
			return l.Type < r.Type
		}
		return true
	})
}

func (o *Override) ApplyTo(s *Schema) error {
	switch o.Type {
	case "model":
		m := s.Models.Get(o.Path.Drop(1), o.Path.Last())
		if m == nil {
			return errors.Errorf("unable to find model at path [%s]", strings.Join(o.Path, "/"))
		}
		return applyModelProperty(m, o.Prop, o.Val)
	case "field":
		mPath := o.Path.Drop(1)
		fKey := o.Path.Last()
		m := s.Models.Get(mPath.Drop(1), mPath.Last())
		if m == nil {
			return errors.Errorf("unable to find model at path [%s]", strings.Join(o.Path, "/"))
		}
		_, f := m.Fields.Get(fKey)
		return applyFieldProperty(m, f, o.Prop, o.Val)
	default:
		return errors.Errorf("unhandled override type [%s]", o.Type)
	}
}

func applyModelProperty(m *model.Model, prop string, val interface{}) error {
	switch prop {
	case KeyTitle:
		m.Title = fmt.Sprintf("%v", val)
	case KeyPlural:
		m.Plural = fmt.Sprintf("%v", val)
	default:
		return errors.Errorf("unhandled model property [%s]", prop)
	}
	return nil
}

func applyFieldProperty(m *model.Model, f *field.Field, prop string, val interface{}) error {
	switch prop {
	case KeyTitle:
		f.Title = fmt.Sprintf("%v", val)
	case KeyPlural:
		f.Plural = fmt.Sprintf("%v", val)
	default:
		return errors.Errorf("unhandled field property [%s]", prop)
	}
	return nil
}
