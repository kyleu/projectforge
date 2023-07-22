//go:build js

package main

import "syscall/js"

func Audit(msg string, code bool, args ...any) {
	js.Global().Call("audit", append([]any{msg, code}, args...)...)
}
