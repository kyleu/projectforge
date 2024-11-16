package types

import "projectforge.dev/projectforge/app/util"

const KeyBit = "bit"

type Bit struct{}

var _ Type = (*Bit)(nil)

func (x *Bit) Key() string {
	return KeyBit
}

func (x *Bit) Sortable() bool {
	return true
}

func (x *Bit) Scalar() bool {
	return true
}

func (x *Bit) String() string {
	return x.Key()
}

func (x *Bit) From(v any) any {
	if x, err := util.ParseBool(v, "", true); err == nil {
		return x
	}
	return invalidInput(x.Key(), v)
}

func (x *Bit) Default(string) any {
	return 0
}

func NewBit() *Wrapped {
	return Wrap(&Bit{})
}
