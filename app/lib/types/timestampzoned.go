// Content managed by Project Forge, see [projectforge.md] for details.
package types

import "time"

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

func (x *TimestampZoned) From(v interface{}) interface{} {
	return invalidInput(x.Key(), x)
}

func (x *TimestampZoned) Default(string) interface{} {
	return time.Now()
}

func NewTimestampZoned() *Wrapped {
	return Wrap(&TimestampZoned{})
}
