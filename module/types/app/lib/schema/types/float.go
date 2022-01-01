package types

const KeyFloat = "float"

type Float struct {
	Bits int
}

var _ Type = (*Float)(nil)

func (x *Float) Key() string {
	return KeyFloat
}

func (x *Float) String() string {
	return x.Key()
}

func (x *Float) Sortable() bool {
	return true
}

func (x *Float) From(v interface{}) interface{} {
	return invalidInput(x.Key(), x)
}

func (x *Float) Default(string) interface{} {
	return 1.0
}

func NewFloat(bits int) *Wrapped {
	return Wrap(&Float{Bits: bits})
}
