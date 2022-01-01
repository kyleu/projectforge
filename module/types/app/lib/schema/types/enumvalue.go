package types

const KeyEnumValue = "enumValue"

type EnumValue struct{}

var _ Type = (*EnumValue)(nil)

func (x *EnumValue) Key() string {
	return KeyEnumValue
}

func (x *EnumValue) String() string {
	return x.Key()
}

func (x *EnumValue) Sortable() bool {
	return true
}

func (x *EnumValue) From(v interface{}) interface{} {
	return invalidInput(x.Key(), x)
}

func (x *EnumValue) Default(key string) interface{} {
	return key
}

func NewEnumValue() *Wrapped {
	return Wrap(&EnumValue{})
}
