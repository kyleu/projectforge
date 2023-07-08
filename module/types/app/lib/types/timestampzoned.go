package types

import "{{{ .Package }}}/app/util"

const KeyTimestampZoned = "timestampZoned"

type TimestampZoned struct{}

var _ Type = (*TimestampZoned)(nil)

func (x *TimestampZoned) Key() string {
	return KeyTimestampZoned
}

func (x *TimestampZoned) Sortable() bool {
	return true
}

func (x *TimestampZoned) Scalar() bool {
	return false
}

func (x *TimestampZoned) String() string {
	return x.Key()
}

func (x *TimestampZoned) From(v any) any {
	return invalidInput(x.Key(), v)
}

func (x *TimestampZoned) Default(string) any {
	return util.TimeCurrent()
}

func NewTimestampZoned() *Wrapped {
	return Wrap(&TimestampZoned{})
}
