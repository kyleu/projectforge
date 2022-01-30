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
	From(v interface{}) interface{}
	Default(key string) interface{}
}

type Types []Type
