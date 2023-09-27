// Package types - Content managed by Project Forge, see [projectforge.md] for details.
package types

const KeyUUID = "uuid"

type UUID struct{}

var _ Type = (*UUID)(nil)

func (x *UUID) Key() string {
	return KeyUUID
}

func (x *UUID) Sortable() bool {
	return true
}

func (x *UUID) Scalar() bool {
	return false
}

func (x *UUID) String() string {
	return x.Key()
}

func (x *UUID) From(v any) any {
	return invalidInput(x.Key(), v)
}

func (x *UUID) Default(string) any {
	return "00000000-0000-0000-0000-000000000000"
}

func NewUUID() *Wrapped {
	return Wrap(&UUID{})
}
