//go:build js
package main

import "syscall/js"

func Audit(msg string, codes ...any) {
	js.Global().Call("audit", append([]any{msg}, codes...)...)
}
