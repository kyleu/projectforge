//go:build js
// +build js

package wasm

import (
	"fmt"
	"syscall/js"

	"{{{ .Package }}}/app/util"
)

func convertArg(arg js.Value) any {
	switch arg.Type() {
	case js.TypeString:
		return arg.String()
	case js.TypeBoolean:
		return arg.Bool()
	case js.TypeNumber:
		return arg.Float()
	case js.TypeObject:
		if arg.Length() > 0 {
			return arrayFrom(arg)
		}
		return valueMapFrom(arg)
	case js.TypeNull:
		return nil
	default:
		return fmt.Sprintf("unhandled type [%T]", arg)
	}
}

func convertArgs(args []js.Value) []any {
	if len(args) == 0 {
		return []any{}
	}
	cleanArgs := make([]any, 0, len(args)-1)
	for _, x := range args[1:] {
		cleanArgs = append(cleanArgs, convertArg(x))
	}
	return cleanArgs
}

func arrayFrom(j js.Value) []any {
	if j.Type() != js.TypeObject {
		return []any{"<error>: input is not a JavaScript array"}
	}
	result := make([]any, 0, j.Length())
	for i := 0; i < j.Length(); i++ {
		result = append(result, convertArg(j.Index(i)))
	}
	return result
}

func valueMapFrom(j js.Value) util.ValueMap {
	if j.Type() != js.TypeObject {
		return util.ValueMap{"<error>": "input is not a JavaScript object"}
	}
	keys := js.Global().Get("Object").Call("keys", j)
	length := keys.Length()
	result := make(util.ValueMap, length)
	for i := 0; i < length; i++ {
		key := keys.Index(i).String()
		result[key] = convertArg(j.Get(key))
	}
	return result
}
