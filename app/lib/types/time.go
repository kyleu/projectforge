// Package types - Content managed by Project Forge, see [projectforge.md] for details.
package types

import "projectforge.dev/projectforge/app/util"

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

func (x *Time) From(v any) any {
	return invalidInput(x.Key(), v)
}

func (x *Time) Default(string) any {
	return util.TimeCurrent().Format("15:04:05")
}

func NewTime() *Wrapped {
	return Wrap(&Time{})
}
