// $PF_GENERATE_ONCE$
//go:build js
// +build js

package wasm

import "syscall/js"

func (w *WASM) wireFunctions() {
	js.Global().Set("log", js.FuncOf(w.logInfo))
	js.Global().Set("sendMessage", js.FuncOf(w.sendMessage))
}

func (w *WASM) logInfo(this js.Value, args []js.Value) any {
	if w.logger == nil {
		println("ERROR: no logger available")
		return nil
	}
	if len(args) == 0 {
		w.logger.Warnf("at least one argument is required")
		return nil
	}
	msg := args[0].String()
	w.logger.Infof(msg, convertArgs(args[1:])...)
	return nil
}

func (w *WASM) sendMessage(this js.Value, args []js.Value) any {
	if len(args) != 1 {
		return "exactly one argument is required"
	}
	arg := args[0].Get("t").String()
	w.logger.Infof("got message of type [%s]", arg)
	return arg
}
