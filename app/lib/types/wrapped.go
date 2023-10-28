// Package types - Content managed by Project Forge, see [projectforge.md] for details.
package types

type Wrapped struct {
	K string `json:"k"`
	T Type   `json:"t,omitempty"`
}

var _ Type = (*Wrapped)(nil)

func Wrap(t Type) *Wrapped {
	w, ok := t.(*Wrapped)
	if ok {
		return w
	}
	return &Wrapped{K: t.Key(), T: t}
}

func (x *Wrapped) Key() string {
	return x.K
}

func (x *Wrapped) Sortable() bool {
	return x.T.Sortable()
}

func (x *Wrapped) Scalar() bool {
	return x.T.Scalar()
}

func (x *Wrapped) String() string {
	return x.T.String()
}

func (x *Wrapped) From(v any) any {
	return x.T.From(v)
}

func (x *Wrapped) Default(key string) any {
	return x.T.Default(key)
}

func (x *Wrapped) Equals(tgt *Wrapped) bool {
	return x.String() == tgt.String()
}

func (x *Wrapped) IsOption() bool {
	_, ok := x.T.(*Option)
	return ok
}

func (x *Wrapped) EnumKey() string {
	e, ok := x.T.(*Enum)
	if !ok {
		return ""
	}
	return e.Ref
}

func (x *Wrapped) ListType() *Wrapped {
	l := TypeAs[*List](x)
	if l != nil {
		return l.V
	}
	return nil
}
