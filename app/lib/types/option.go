// Package types - Content managed by Project Forge, see [projectforge.md] for details.
package types

const KeyOption = "option"

type Option struct {
	V *Wrapped `json:"v"`
}

var _ Type = (*Option)(nil)

func (x *Option) Key() string {
	return KeyOption
}

func (x *Option) Sortable() bool {
	return x.V.Sortable()
}

func (x *Option) Scalar() bool {
	return false
}

func (x *Option) String() string {
	return "*" + x.V.String()
}

func (x *Option) From(v any) any {
	return x.V.From(v)
}

func NewOption(t *Wrapped) *Wrapped {
	return Wrap(&Option{V: t})
}

func (x *Option) Default(string) any {
	return "∅"
}
