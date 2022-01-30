// Content managed by Project Forge, see [projectforge.md] for details.
package types

const KeyChar = "char"

type Char struct{}

var _ Type = (*Char)(nil)

func (x *Char) Key() string {
	return KeyChar
}

func (x *Char) Sortable() bool {
	return true
}

func (x *Char) Scalar() bool {
	return true
}

func (x *Char) String() string {
	return x.Key()
}

func (x *Char) From(v interface{}) interface{} {
	return invalidInput(x.Key(), x)
}

func (x *Char) Default(string) interface{} {
	return 'x'
}

func NewChar() *Wrapped {
	return Wrap(&Char{})
}
