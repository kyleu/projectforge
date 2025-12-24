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

func (u *Union) Key() string {
	return KeyUnion
}

func (u *Union) Sortable() bool {
	return lo.EveryBy(u.V, func(v *Wrapped) bool {
		return v.Sortable()
	})
}

func (u *Union) Scalar() bool {
	return false
}

func (u *Union) String() string {
	ss := util.NewStringSliceWithSize(len(u.V))
	for _, v := range u.V {
		ss.Push(v.String())
	}
	return fmt.Sprintf("[%s]", ss.JoinCommas())
}

func (u *Union) From(v any) any {
	return invalidInput(u.Key(), v)
}

func (u *Union) Default(x string) any {
	if len(u.V) == 0 {
		return KeyNilString
	}
	return u.V[0].Default(x)
}

func NewUnion(t ...*Wrapped) *Wrapped {
	return Wrap(&Union{V: t})
}
