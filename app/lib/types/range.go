// Package types - Content managed by Project Forge, see [projectforge.md] for details.
package types

import (
	"fmt"
)

const KeyRange = "range"

var _ Type = (*Range)(nil)

type Range struct {
	V *Wrapped `json:"v"`
}

func (x *Range) Key() string {
	return KeyRange
}

func (x *Range) Sortable() bool {
	return x.V.Sortable()
}

func (x *Range) Scalar() bool {
	return false
}

func (x *Range) String() string {
	return fmt.Sprintf("range[%s]", x.V.String())
}

func (x *Range) From(v any) any {
	return invalidInput(x.Key(), v)
}

func (x *Range) Default(string) any {
	return ""
}

func NewRange(t *Wrapped) *Wrapped {
	return Wrap(&Range{V: t})
}
