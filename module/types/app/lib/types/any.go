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

func (x *Any) From(v interface{}) interface{} {
	return v
}

func (x *Any) Default(string) interface{} {
	return nil
}

func NewAny() *Wrapped {
	return Wrap(&Any{})
}
