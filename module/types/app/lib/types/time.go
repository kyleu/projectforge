package types

import "time"

const KeyTime = "time"

type Time struct{}

var _ Type = (*Time)(nil)

func (x *Time) Key() string {
	return KeyTime
}

func (x *Time) Sortable() bool {
	return true
}

func (x *Time) Scalar() bool {
	return false
}

func (x *Time) String() string {
	return x.Key()
}

func (x *Time) From(v interface{}) interface{} {
	return invalidInput(x.Key(), x)
}

func (x *Time) Default(string) interface{} {
	return time.Now().Format("15:04:05")
}

func NewTime() *Wrapped {
	return Wrap(&Time{})
}
