//go:build js
// +build js

package main

import "{{{ .Package }}}/app/wasm"

func main() {
	w, err := wasm.NewWASM()
	if err != nil {
		panic(err)
	}
	err = w.Run()
	if err != nil {
		panic(err)
	}
}
