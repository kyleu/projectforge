package types

import "projectforge.dev/projectforge/app/util"

type Wrapped struct {
	K string `json:"k"`
	T Type   `json:"t,omitempty"`
}

var _ Type = (*Wrapped)(nil)

func Wrap(t Type) *Wrapped {
	if w, err := util.Cast[*Wrapped](t); err == nil {
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
	_, err := util.Cast[*Option](x.T)
	return err == nil
}

func (x *Wrapped) EnumKey() string {
	e, err := util.Cast[*Enum](x.T)
	if err != nil {
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

func (x *Wrapped) OrderedMapType() *Wrapped {
	l := TypeAs[*OrderedMap](x)
	if l != nil {
		return l.V
	}
	return nil
}
