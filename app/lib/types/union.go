package types

import (
	"fmt"

	"github.com/samber/lo"

	"projectforge.dev/projectforge/app/util"
)

const KeyUnion = "union"

var _ Type = (*Union)(nil)

type Union struct {
	V []*Wrapped `json:"v"`
}

func (x *Union) Key() string {
	return KeyUnion
}

func (x *Union) Sortable() bool {
	return lo.EveryBy(x.V, func(v *Wrapped) bool {
		return v.Sortable()
	})
}

func (x *Union) Scalar() bool {
	return false
}

func (x *Union) String() string {
	ss := util.NewStringSliceWithSize(len(x.V))
	for _, v := range x.V {
		ss.Push(v.String())
	}
	return fmt.Sprintf("[%s]", ss.JoinCommas())
}

func (x *Union) From(v any) any {
	return invalidInput(x.Key(), v)
}

func (x *Union) Default(string) any {
	if len(x.V) == 0 {
		return "<nil>"
	}
	return x.V[0].Default("")
}

func NewUnion(t ...*Wrapped) *Wrapped {
	return Wrap(&Union{V: t})
}
