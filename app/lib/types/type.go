// Package types - Content managed by Project Forge, see [projectforge.md] for details.
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

func TypeAs[T Type](t Type) T {
	l, ok := t.(T)
	if ok {
		return l
	}
	w, ok := t.(*Wrapped)
	if !ok {
		var ret T
		return ret
	}
	l, ok = w.T.(T)
	if !ok {
		var ret T
		return ret
	}
	return l
}

func IsString(t Type) bool {
	return t.Key() == KeyString
}

func IsBool(t Type) bool {
	return t.Key() == KeyBool
}

func IsInt(t Type) bool {
	return t.Key() == KeyInt
}

func IsList(t Type) bool {
	return t.Key() == KeyList
}

func IsStringList(t Type) bool {
	l := TypeAs[*List](t)
	if l == nil {
		return false
	}
	return l.V.T.Key() == KeyString
}
