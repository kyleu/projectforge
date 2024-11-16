package types

import "{{{ .Package }}}/app/util"

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

func (x *Char) From(v any) any {
	if x, err := util.ParseString(v, "", true); err == nil && len(x) > 0 {
		return x[0]
	}
	return invalidInput(x.Key(), v)
}

func (x *Char) Default(string) any {
	return 'x'
}

func NewChar() *Wrapped {
	return Wrap(&Char{})
}
