// $PF_GENERATE_ONCE$
//go:build js
// +build js

package wasm

import (
	"syscall/js"
	"time"

	"{{{ .Package }}}/app/util"
)

func (w *WASM) Testbed(args ...any) js.Value {
	return w.call("testbed", args...)
}

func (w *WASM) OnMessage(typ string, msg any) js.Value {
	return w.call("onMessage", typ, msg)
}

func (w *WASM) Log(level string, occurred time.Time, loggerName string, message string, caller util.ValueMap, stack string, fields util.ValueMap) js.Value {
	m := util.ValueMap{"level": level, "message": message, "caller": caller.AsMap(), "occurred": util.TimeToJS(&occurred)}
	if stack != "" {
		m["stack"] = stack
	}
	if len(fields) > 0 {
		m["fields"] = fields.AsMap()
	}
	return w.OnMessage("log", m.AsMap())
}

func (w *WASM) call(fn string, args ...any) js.Value {
	if x := js.Global().Get(fn); x.IsUndefined() {
		w.logger.Warnf("function [%s] is not defined", fn)
		return js.Undefined()
	} else {
		return js.Global().Call(fn, args...)
	}
}
