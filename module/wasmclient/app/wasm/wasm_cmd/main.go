//go:build js
// +build js

package main

import "{{{ .Package }}}/app/wasm"

func main() {
	w := wasm.NewWASM()
	<-w.CloseCh
}
