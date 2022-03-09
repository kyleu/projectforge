// Content managed by Project Forge, see [projectforge.md] for details.
package types

import (
	"fmt"

	"projectforge.dev/projectforge/app/util"
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
	return false
}

func (x *List) String() string {
	return fmt.Sprintf("[]%s", x.V.String())
}

func (x *List) From(v interface{}) interface{} {
	switch t := v.(type) {
	case string:
		lt := util.StringSplitAndTrim(t, ",")
		return lt
	default:
		return invalidInput(x.Key(), t)
	}
}

func (x *List) Default(string) interface{} {
	return emptyList
}

func NewList(t *Wrapped) *Wrapped {
	return Wrap(&List{V: t})
}

func NewListSized(t *Wrapped, size int) *Wrapped {
	return Wrap(&List{V: t, Size: size})
}
