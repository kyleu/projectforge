package types

import (
	"fmt"

	"{{{ .Package }}}/app/util"
)

const KeyMap = "map"

var _ Type = (*Map)(nil)

type Map struct {
	K *Wrapped `json:"k"`
	V *Wrapped `json:"v"`
}

func (x *Map) Key() string {
	return KeyMap
}

func (x *Map) Sortable() bool {
	return x.K.Sortable()
}

func (x *Map) Scalar() bool {
	return true
}

func (x *Map) String() string {
	k, v := x.K.String(), x.V.String()
	if k == KeyString && v == KeyAny {
		return x.Key()
	}
	return fmt.Sprintf("%s[%s]%s", x.Key(), k, v)
}

func (x *Map) From(v any) any {
	if x, err := util.ParseMap(v, "", true); err == nil {
		return x
	}
	return invalidInput(x.Key(), v)
}

func (x *Map) Default(string) any {
	return emptyMap
}

func NewMap(k *Wrapped, v *Wrapped) *Wrapped {
	return Wrap(&Map{K: k, V: v})
}

func NewStringKeyedMap() *Wrapped {
	return NewMap(NewString(), NewAny())
}
