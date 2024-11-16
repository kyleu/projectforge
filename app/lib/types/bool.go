package types

import "projectforge.dev/projectforge/app/util"

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
	if x, err := util.ParseBool(v, "", true); err == nil {
		return x
	}
	return invalidInput(x.Key(), v)
}

func (x *Bool) Default(string) any {
	return false
}

func NewBool() *Wrapped {
	return Wrap(&Bool{})
}
