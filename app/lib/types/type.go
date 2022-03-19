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
