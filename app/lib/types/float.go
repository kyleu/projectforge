// Package types - Content managed by Project Forge, see [projectforge.md] for details.
package types

const KeyFloat = "float"

type Float struct {
	Bits int `json:"bits,omitempty"`
}

var _ Type = (*Float)(nil)

func (x *Float) Key() string {
	return KeyFloat
}

func (x *Float) Sortable() bool {
	return true
}

func (x *Float) Scalar() bool {
	return true
}

func (x *Float) String() string {
	return x.Key()
}

func (x *Float) From(v any) any {
	return invalidInput(x.Key(), v)
}

func (x *Float) Default(string) any {
	return 1.0
}

func NewFloat(bits int) *Wrapped {
	return Wrap(&Float{Bits: bits})
}
