package types

import (
	"fmt"

	"projectforge.dev/projectforge/app/util"
)

const KeyString = "string"

type String struct {
	MinLength int    `json:"minLength,omitzero"`
	MaxLength int    `json:"maxLength,omitzero"`
	Pattern   string `json:"pattern,omitzero"`
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
	if x, err := util.ParseString(v, "", true); err == nil {
		return x
	}
	return invalidInput(x.Key(), v)
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
