package types

import (
	"fmt"

	"{{{ .Package }}}/app/util"
)

const KeyList = "list"

var _ Type = (*List)(nil)

type List struct {
	V    *Wrapped `json:"v"`
	Size int      `json:"size,omitempty"`
}

func (x *List) Key() string {
	return KeyList
}

func (x *List) Sortable() bool {
	return x.V.Sortable()
}

func (x *List) Scalar() bool {
	return true
}

func (x *List) String() string {
	return fmt.Sprintf("[]%s", x.V.String())
}

func (x *List) From(v any) any {
	if x, err := util.ParseArray(v, "", true, true); err == nil {
		return x
	}
	return invalidInput(x.Key(), v)
}

func (x *List) Default(string) any {
	return emptyList
}

func NewList(t *Wrapped) *Wrapped {
	return Wrap(&List{V: t})
}

func NewListSized(t *Wrapped, size int) *Wrapped {
	return Wrap(&List{V: t, Size: size})
}
