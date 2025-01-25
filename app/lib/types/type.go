package types

import (
	"fmt"

	"projectforge.dev/projectforge/app/util"
)

type Type interface {
	fmt.Stringer
	Key() string
	Sortable() bool
	Scalar() bool
	From(v any) any
	Default(key string) any
}

type Types []Type

func TypeAs[T Type](t Type) T {
	l, err := util.Cast[T](t)
	if err == nil {
		return l
	}
	w, err := util.Cast[*Wrapped](t)
	if err != nil {
		var ret T
		return ret
	}
	l, err = util.Cast[T](w.T)
	if err != nil {
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
