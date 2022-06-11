// Content managed by Project Forge, see [projectforge.md] for details.
package types

import (
	"fmt"
)

type Type interface {
	Key() string
	Sortable() bool
	Scalar() bool
	fmt.Stringer
	From(v any) any
	Default(key string) any
}

type Types []Type

func TypeAs[T any](t Type) T {
	l, ok := t.(T)
	if !ok {
		w, ok := t.(*Wrapped)
		if ok {
			l, ok = w.T.(T)
		}
		if !ok {
			var ret T
			return ret
		}
	}
	return l
}
