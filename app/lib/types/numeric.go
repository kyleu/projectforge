package types

const KeyNumeric = "bignum"

type Numeric struct{}

var _ Type = (*Numeric)(nil)

func (x *Numeric) Key() string {
	return KeyNumeric
}

func (x *Numeric) Sortable() bool {
	return true
}

func (x *Numeric) Scalar() bool {
	return true
}

func (x *Numeric) String() string {
	return x.Key()
}

func (x *Numeric) From(v any) any {
	return invalidInput(x.Key(), v)
}

func (x *Numeric) Default(string) any {
	return 1.0
}

func NewNumeric() *Wrapped {
	return Wrap(&Numeric{})
}
