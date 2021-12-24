package field

import (
	"reflect"

	"github.com/kyleu/projectforge/app/schema/types"
)

func NewFieldByType(key string, t reflect.Type, ro bool, md *Metadata) *Field {
	return &Field{Key: key, Type: fromReflect(t), ReadOnly: ro, Metadata: md}
}

// nolint
func fromReflect(t reflect.Type) *types.Wrapped {
	switch t.Kind() {
	case reflect.Invalid:
		return types.NewError("can't reflect invalid")
	case reflect.Bool:
		return types.NewBool()
	case reflect.Int:
		return types.NewInt(0)
	case reflect.Int8:
		return types.NewInt(8)
	case reflect.Int16:
		return types.NewInt(16)
	case reflect.Int32:
		return types.NewInt(32)
	case reflect.Int64:
		return types.NewInt(64)
	case reflect.Uint:
		return types.NewUnsignedInt(0)
	case reflect.Uint8:
		return types.NewUnsignedInt(8)
	case reflect.Uint16:
		return types.NewUnsignedInt(16)
	case reflect.Uint32:
		return types.NewUnsignedInt(32)
	case reflect.Uint64:
		return types.NewUnsignedInt(64)
	case reflect.Uintptr:
		return types.NewError("can't reflect uint ponters")
	case reflect.Float32:
		return types.NewFloat(32)
	case reflect.Float64:
		return types.NewFloat(64)
	case reflect.Complex64:
		return types.NewError("can't reflect complex")
	case reflect.Complex128:
		return types.NewError("can't reflect complex")
	case reflect.Array:
		return types.NewList(fromReflect(t.Elem()))
	case reflect.Chan:
		return types.NewError("can't reflect channels")
	case reflect.Func:
		return types.NewError("can't reflect functions")
	case reflect.Interface:
		return types.NewError("can't reflect interfaces")
	case reflect.Map:
		return types.NewMap(fromReflect(t.Key()), fromReflect(t.Elem()))
	case reflect.Ptr:
		return types.NewOption(fromReflect(t.Elem()))
	case reflect.Slice:
		return types.NewList(fromReflect(t.Elem()))
	case reflect.String:
		return types.NewString()
	case reflect.Struct:
		return types.NewError("can't reflect structs")
	case reflect.UnsafePointer:
		return types.NewError("can't reflect unsafe pointers")
	default:
		return types.NewUnknown(t.Kind().String())
	}
}
