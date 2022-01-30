// Content managed by Project Forge, see [projectforge.md] for details.
package types

import (
	"fmt"
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
	return fmt.Sprintf("%s[%s]%s", x.Key(), x.K.String(), x.V.String())
}

func (x *Map) From(v interface{}) interface{} {
	return invalidInput(x.Key(), x)
}

func (x *Map) Default(string) interface{} {
	return ""
}

func NewMap(k *Wrapped, v *Wrapped) *Wrapped {
	return Wrap(&Map{K: k, V: v})
}
