package types

import (
	"fmt"
)

type Type interface {
	Key() string
	fmt.Stringer
	Sortable() bool
	From(v interface{}) interface{}
	Default(key string) interface{}
}
