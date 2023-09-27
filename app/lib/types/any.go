// Package types - Content managed by Project Forge, see [projectforge.md] for details.
package types

const KeyAny = "any"

type Any struct{}

var _ Type = (*Any)(nil)

func (x *Any) Key() string {
	return KeyAny
}

func (x *Any) Sortable() bool {
	return false
}

func (x *Any) Scalar() bool {
	return false
}

func (x *Any) String() string {
	return x.Key()
}

func (x *Any) From(v any) any {
	return v
}

func (x *Any) Default(string) any {
	return nil
}

func NewAny() *Wrapped {
	return Wrap(&Any{})
}
