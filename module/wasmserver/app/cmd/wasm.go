//go:build !js
// +build !js

package cmd

import (
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

const keyWASM = "wasm"

func wasmCmd() *cobra.Command {
	short := "Starts the server and exposes a WebAssembly application to scripts"
	f := func(*cobra.Command, []string) error { return startWASM(_flags) }
	ret := newCmd(keyWASM, short, f)
	return ret
}

func startWASM(_ *Flags) error {
	return errors.New("The WASM command can only be run from WebAssembly")
}
