// Package types - Content managed by Project Forge, see [projectforge.md] for details.
package types

import (
	"fmt"
)

const KeySet = "set"

type Set struct {
	V *Wrapped `json:"v"`
}

var _ Type = (*Set)(nil)

func (x *Set) Key() string {
	return KeySet
}

func (x *Set) Sortable() bool {
	return x.V.Sortable()
}

func (x *Set) Scalar() bool {
	return false
}

func (x *Set) String() string {
	return fmt.Sprintf("%s[%s]", x.Key(), x.V.String())
}

func (x *Set) From(v any) any {
	return invalidInput(x.Key(), v)
}

func (x *Set) Default(string) any {
	return emptyList
}

func NewSet() *Wrapped {
	return Wrap(&Set{})
}
