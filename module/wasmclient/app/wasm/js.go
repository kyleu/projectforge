//go:build js
// +build js

package main

import (
	"syscall/js"
	"time"

	"{{{ .Package }}}/app/util"
)

func Testbed(args ...any) js.Value {
	return call("testbed", args...)
}

func OnMessage(typ string, msg any) js.Value {
	return call("onMessage", typ, msg)
}

func Log(level string, occurred time.Time, loggerName string, message string, caller util.ValueMap, stack string, fields util.ValueMap) js.Value {
	m := util.ValueMap{"level": level, "message": message, "caller": caller.AsMap(), "occurred": util.TimeToJS(&occurred)}
	if stack != "" {
		m["stack"] = stack
	}
	if len(fields) > 0 {
		m["fields"] = fields.AsMap()
	}
	return OnMessage("log", m.AsMap())
}

func call(fn string, args ...any) js.Value {
	if x := js.Global().Get(fn); x.IsUndefined() {
		_rootLogger.Warnf("function [%s] is not defined", fn)
		return js.Undefined()
	} else {
		return js.Global().Call(fn, args...)
	}
}
