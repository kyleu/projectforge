// Package types - Content managed by Project Forge, see [projectforge.md] for details.
package types

import (
	"reflect"
)

func FromReflect(t reflect.Type) *Wrapped {
	switch t.Kind() {
	case reflect.Invalid:
		return NewError("can't reflect invalid")
	case reflect.Bool:
		return NewBool()
	case reflect.Int:
		return NewInt(0)
	case reflect.Int8:
		return NewInt(8)
	case reflect.Int16:
		return NewInt(16)
	case reflect.Int32:
		return NewInt(32)
	case reflect.Int64:
		return NewInt(64)
	case reflect.Uint:
		return NewUnsignedInt(0)
	case reflect.Uint8:
		return NewUnsignedInt(8)
	case reflect.Uint16:
		return NewUnsignedInt(16)
	case reflect.Uint32:
		return NewUnsignedInt(32)
	case reflect.Uint64:
		return NewUnsignedInt(64)
	case reflect.Uintptr:
		return NewError("can't reflect uint ponters")
	case reflect.Float32:
		return NewFloat(32)
	case reflect.Float64:
		return NewFloat(64)
	case reflect.Complex64:
		return NewError("can't reflect complex")
	case reflect.Complex128:
		return NewError("can't reflect complex")
	case reflect.Array:
		return NewList(FromReflect(t.Elem()))
	case reflect.Chan:
		return NewError("can't reflect channels")
	case reflect.Func:
		return NewError("can't reflect functions")
	case reflect.Interface:
		return NewAny()
	case reflect.Map:
		return NewMap(FromReflect(t.Key()), FromReflect(t.Elem()))
	case reflect.Ptr:
		return NewOption(FromReflect(t.Elem()))
	case reflect.Slice:
		return NewList(FromReflect(t.Elem()))
	case reflect.String:
		return NewString()
	case reflect.Struct:
		return NewError("can't reflect structs")
	case reflect.UnsafePointer:
		return NewError("can't reflect unsafe pointers")
	default:
		return NewUnknown(t.Kind().String())
	}
}
