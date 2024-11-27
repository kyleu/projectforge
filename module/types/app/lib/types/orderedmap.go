package types

import (
	"fmt"

	"{{{ .Package }}}/app/util"
)

const KeyOrderedMap = "orderedMap"

var _ Type = (*OrderedMap)(nil)

type OrderedMap struct {
	V *Wrapped `json:"v"`
}

func (x *OrderedMap) Key() string {
	return KeyOrderedMap
}

func (x *OrderedMap) Sortable() bool {
	return x.V.Sortable()
}

func (x *OrderedMap) Scalar() bool {
	return true
}

func (x *OrderedMap) String() string {
	return fmt.Sprintf("[]%s", x.V.String())
}

func (x *OrderedMap) From(v any) any {
	if x, err := util.ParseOrderedMap(v, "", true); err == nil {
		return x
	}
	return invalidInput(x.Key(), v)
}

func (x *OrderedMap) Default(string) any {
	return emptyMap
}

func NewOrderedMap(t *Wrapped) *Wrapped {
	return Wrap(&OrderedMap{V: t})
}
