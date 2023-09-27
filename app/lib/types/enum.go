// Package types - Content managed by Project Forge, see [projectforge.md] for details.
package types

import (
	"fmt"
)

const KeyEnum = "enum"

type Enum struct {
	Ref string `json:"ref"`
}

var _ Type = (*Enum)(nil)

func (x *Enum) Key() string {
	return KeyEnum
}

func (x *Enum) Sortable() bool {
	return true
}

func (x *Enum) Scalar() bool {
	return true
}

func (x *Enum) String() string {
	return fmt.Sprintf("%s(%s)", x.Key(), x.Ref)
}

func (x *Enum) From(v any) any {
	switch t := v.(type) {
	case string:
		return t
	default:
		return invalidInput(x.Key(), t)
	}
}

func (x *Enum) Default(_ string) any {
	return ""
}

func NewEnum(ref string) *Wrapped {
	return Wrap(&Enum{Ref: ref})
}
