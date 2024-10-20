//go:build js
// +build js

package wasm

import (
	"fmt"
	"syscall/js"
)

func convertArgs(args []js.Value) []any {
	cleanArgs := make([]any, 0, len(args)-1)
	for _, x := range args[1:] {
		switch x.Type() {
		case js.TypeString:
			cleanArgs = append(cleanArgs, x.String())
		case js.TypeBoolean:
			cleanArgs = append(cleanArgs, x.Bool())
		case js.TypeNumber:
			cleanArgs = append(cleanArgs, x.Float())
		case js.TypeNull:
			cleanArgs = append(cleanArgs, "<null>")
		default:
			cleanArgs = append(cleanArgs, fmt.Sprintf("unhandled type [%T]", x))
		}
	}
	return cleanArgs
}
