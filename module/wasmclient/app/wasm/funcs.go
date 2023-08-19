// $PF_GENERATE_ONCE$
//go:build js

package main

import (
	"fmt"
	"syscall/js"
)

func wireFunctions() {
	js.Global().Set("increment", js.FuncOf(increment))
	js.Global().Set("log", js.FuncOf(log))
	js.Global().Set("send", js.FuncOf(send))
}

func increment(this js.Value, args []js.Value) any {
	if len(args) == 0 {
		return "exactly one argument is required"
	}
	arg := args[0].Int()
	return arg + 1
}

func send(this js.Value, args []js.Value) any {
	if len(args) == 0 {
		return "at least one argument is required"
	}
	arg := args[0].String()
	return arg
}

func log(this js.Value, args []js.Value) any {
	if _rootLogger == nil {
		println("ERROR: no logger available")
		return nil
	}
	if len(args) == 0 {
		_rootLogger.Warnf("at least one argument is required")
		return nil
	}
	msg := args[0].String()
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
	_rootLogger.Infof(msg, cleanArgs...)
	return nil
}
