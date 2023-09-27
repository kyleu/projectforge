// Package types - Content managed by Project Forge, see [projectforge.md] for details.
package types

const KeyUnknown = "unknown"

type Unknown struct {
	X string `json:"x"`
}

var _ Type = (*Unknown)(nil)

func (x *Unknown) Key() string {
	return KeyUnknown
}

func (x *Unknown) Sortable() bool {
	return false
}

func (x *Unknown) Scalar() bool {
	return false
}

func (x *Unknown) String() string {
	return x.Key() + "(" + x.X + ")"
}

func (x *Unknown) From(v any) any {
	return v
}

func (x *Unknown) Default(string) any {
	return x.X
}

func NewUnknown(x string) *Wrapped {
	return Wrap(&Unknown{X: x})
}
