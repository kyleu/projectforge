// Package types - Content managed by Project Forge, see [projectforge.md] for details.
package types

import (
	"fmt"
	"strconv"
)

const KeyInt = "int"

type Int struct {
	Bits     int  `json:"bits,omitempty"`
	Unsigned bool `json:"unsigned,omitempty"`
}

var _ Type = (*Int)(nil)

func (x *Int) Key() string {
	return KeyInt
}

func (x *Int) Sortable() bool {
	return true
}

func (x *Int) Scalar() bool {
	return true
}

func (x *Int) String() string {
	if x.Bits > 0 {
		return fmt.Sprintf("%s%d", x.Key(), x.Bits)
	}
	return x.Key()
}

func (x *Int) From(v any) any {
	switch t := v.(type) {
	case string:
		ret, _ := strconv.ParseInt(t, 10, 32)
		return int(ret)
	case int:
		return t
	case int32:
		return t
	case int64:
		return t
	default:
		return invalidInput(x.Key(), t)
	}
}

func (x *Int) Default(string) any {
	return 0
}

func NewInt(bits int) *Wrapped {
	return Wrap(&Int{Bits: bits})
}

func NewUnsignedInt(bits int) *Wrapped {
	return Wrap(&Int{Bits: bits, Unsigned: true})
}
