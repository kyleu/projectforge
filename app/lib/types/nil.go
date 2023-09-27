// Package types - Content managed by Project Forge, see [projectforge.md] for details.
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

func (x *Nil) From(v any) any {
	switch v {
	case nil:
		return nil
	default:
		return invalidInput(x.Key(), v)
	}
}

func (x *Nil) Default(string) any {
	return "<nil>"
}
