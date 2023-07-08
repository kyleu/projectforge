package types

import "{{{ .Package }}}/app/util"

const KeyTimestamp = "timestamp"

type Timestamp struct{}

var _ Type = (*Timestamp)(nil)

func (x *Timestamp) Key() string {
	return KeyTimestamp
}

func (x *Timestamp) Sortable() bool {
	return true
}

func (x *Timestamp) Scalar() bool {
	return false
}

func (x *Timestamp) String() string {
	return x.Key()
}

func (x *Timestamp) From(v any) any {
	return invalidInput(x.Key(), v)
}

func (x *Timestamp) Default(string) any {
	return util.TimeCurrent()
}

func NewTimestamp() *Wrapped {
	return Wrap(&Timestamp{})
}
