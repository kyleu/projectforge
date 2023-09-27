// Package types - Content managed by Project Forge, see [projectforge.md] for details.
package types

import (
	"fmt"
)

const KeyString = "string"

type String struct {
	MinLength int    `json:"minLength,omitempty"`
	MaxLength int    `json:"maxLength,omitempty"`
	Pattern   string `json:"pattern,omitempty"`
}

var _ Type = (*String)(nil)

func (x *String) Key() string {
	return KeyString
}

func (x *String) Sortable() bool {
	return true
}

func (x *String) Scalar() bool {
	return true
}

func (x *String) String() string {
	if x.MaxLength > 0 {
		return fmt.Sprintf("%s(%d)", x.Key(), x.MaxLength)
	}
	return x.Key()
}

func (x *String) From(v any) any {
	switch t := v.(type) {
	case string:
		return t
	default:
		return invalidInput(x.Key(), t)
	}
}

func (x *String) Default(string) any {
	return ""
}

func NewString() *Wrapped {
	return Wrap(&String{})
}

func NewStringArgs(minLength int, maxLength int, pattern string) *Wrapped {
	return Wrap(&String{MinLength: minLength, MaxLength: maxLength, Pattern: pattern})
}
