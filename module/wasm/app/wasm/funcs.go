// $PF_IGNORE$
package main

import (
	"syscall/js"
)

func wireFunctions() {
	js.FuncOf(foo)
}

func foo(this js.Value, args []js.Value) any {
	if len(args) != 1 {
		return "Invalid no of arguments passed"
	}
	arg := args[0].String()
	return arg
}
