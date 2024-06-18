package types

import "projectforge.dev/projectforge/app/util"

const KeyDate = "date"

type Date struct{}

var _ Type = (*Date)(nil)

func (x *Date) Key() string {
	return KeyDate
}

func (x *Date) Sortable() bool {
	return true
}

func (x *Date) Scalar() bool {
	return false
}

func (x *Date) String() string {
	return x.Key()
}

func (x *Date) From(v any) any {
	return invalidInput(x.Key(), v)
}

func (x *Date) Default(string) any {
	return util.TimeCurrent()
}

func NewDate() *Wrapped {
	return Wrap(&Date{})
}
