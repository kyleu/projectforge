package types

import (
	"strings"

	"{{{ .Package }}}/app/util"
)

const KeyBool = "bool"

type Bool struct{}

var _ Type = (*Bool)(nil)

func (x *Bool) Key() string {
	return KeyBool
}

func (x *Bool) Sortable() bool {
	return true
}

func (x *Bool) Scalar() bool {
	return true
}

func (x *Bool) String() string {
	return x.Key()
}

func (x *Bool) From(v any) any {
	switch t := v.(type) {
	case bool:
		return t
	case string:
		lt := strings.ToLower(t)
		return lt == util.BoolTrue || lt == "yes" || lt == "t"
	default:
		return invalidInput(x.Key(), t)
	}
}

func (x *Bool) Default(string) any {
	return false
}

func NewBool() *Wrapped {
	return Wrap(&Bool{})
}
