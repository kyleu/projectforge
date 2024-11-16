package types

import "{{{ .Package }}}/app/util"

const KeyEnumValue = "enumValue"

type EnumValue struct{}

var _ Type = (*EnumValue)(nil)

func (x *EnumValue) Key() string {
	return KeyEnumValue
}

func (x *EnumValue) Sortable() bool {
	return true
}

func (x *EnumValue) Scalar() bool {
	return true
}

func (x *EnumValue) String() string {
	return x.Key()
}

func (x *EnumValue) From(v any) any {
	if x, err := util.ParseString(v, "", true); err == nil {
		return x
	}
	return invalidInput(x.Key(), v)
}

func (x *EnumValue) Default(key string) any {
	return key
}

func NewEnumValue() *Wrapped {
	return Wrap(&EnumValue{})
}
