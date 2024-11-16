package types

import (
	"strconv"

	"{{{ .Package }}}/app/util"
)

const KeyFloat = "float"

type Float struct {
	Bits int `json:"bits,omitempty"`
}

var _ Type = (*Float)(nil)

func (x *Float) Key() string {
	return KeyFloat
}

func (x *Float) Sortable() bool {
	return true
}

func (x *Float) Scalar() bool {
	return true
}

func (x *Float) String() string {
	return x.Key()
}

func (x *Float) From(v any) any {
	switch t := v.(type) {
	case string:
		ret, _ := strconv.ParseFloat(t, util.Choose(x.Bits == 0, 64, x.Bits))
		return ret
	case float32:
		return float64(t)
	case float64:
		return t
	case int:
		return float64(t)
	case int32:
		return float64(t)
	case int64:
		return float64(t)
	default:
		return invalidInput(x.Key(), t)
	}
}

func (x *Float) Default(string) any {
	return 1.0
}

func NewFloat(bits int) *Wrapped {
	return Wrap(&Float{Bits: bits})
}
