// Package types - Content managed by Project Forge, see [projectforge.md] for details.
package types

const KeyByte = "byte"

type Byte struct{}

var _ Type = (*Byte)(nil)

func (x *Byte) Key() string {
	return KeyByte
}

func (x *Byte) Sortable() bool {
	return true
}

func (x *Byte) Scalar() bool {
	return true
}

func (x *Byte) String() string {
	return x.Key()
}

func (x *Byte) From(v any) any {
	return invalidInput(x.Key(), v)
}

func (x *Byte) Default(string) any {
	return 0
}

func NewByte() *Wrapped {
	return Wrap(&Byte{})
}
