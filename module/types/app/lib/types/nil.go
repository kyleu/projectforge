package types

const KeyNil = "nil"

type Nil struct{}

var _ Type = (*Nil)(nil)

func (x *Nil) Key() string {
	return KeyNil
}

func (x *Nil) Sortable() bool {
	return false
}

func (x *Nil) Scalar() bool {
	return false
}

func (x *Nil) String() string {
	return x.Key()
}

func (x *Nil) From(v interface{}) interface{} {
	switch v {
	case nil:
		return nil
	default:
		return invalidInput(x.Key(), v)
	}
}

func (x *Nil) Default(string) interface{} {
	return "<nil>"
}
