package types

import (
	"fmt"
	"strconv"

	"{{{ .Package }}}/app/util"
)

const KeyInt = "int"

type Int struct {
	Bits     int  `json:"bits,omitzero"`
	Unsigned bool `json:"unsigned,omitzero"`
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
		ret, _ := strconv.ParseInt(t, 10, util.Choose(x.Bits == 0, 64, x.Bits))
		return int(ret)
	case float32:
		return int(t)
	case float64:
		return int(t)
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
