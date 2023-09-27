// Package types - Content managed by Project Forge, see [projectforge.md] for details.
package types

const KeyXML = "xml"

type XML struct{}

var _ Type = (*XML)(nil)

func (x *XML) Key() string {
	return KeyXML
}

func (x *XML) Sortable() bool {
	return true
}

func (x *XML) Scalar() bool {
	return false
}

func (x *XML) String() string {
	return x.Key()
}

func (x *XML) From(v any) any {
	return invalidInput(x.Key(), v)
}

func (x *XML) Default(string) any {
	return "<element />"
}

func NewXML() *Wrapped {
	return Wrap(&XML{})
}
