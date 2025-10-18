package types

const KeyNumericMap = "numericMap"

type NumericMap struct{}

var _ Type = (*NumericMap)(nil)

func (x *NumericMap) Key() string {
	return KeyNumericMap
}

func (x *NumericMap) Sortable() bool {
	return true
}

func (x *NumericMap) Scalar() bool {
	return false
}

func (x *NumericMap) String() string {
	return x.Key()
}

func (x *NumericMap) From(v any) any {
	return invalidInput(x.Key(), v)
}

func (x *NumericMap) Default(key string) any {
	return emptyMap
}

func NewNumericMap() *Wrapped {
	return Wrap(&NumericMap{})
}
