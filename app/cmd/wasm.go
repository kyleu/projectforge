//go:build !js

// Package cmd - Content managed by Project Forge, see [projectforge.md] for details.
package cmd

import (
	"github.com/muesli/coral"
	"github.com/pkg/errors"
)

const keyWASM = "wasm"

func wasmCmd() *coral.Command {
	short := "Starts the server and exposes a WebAssembly application to scripts"
	f := func(*coral.Command, []string) error { return startWASM(_flags) }
	ret := &coral.Command{Use: keyWASM, Short: short, RunE: f}
	return ret
}

func startWASM(_ *Flags) error {
	return errors.New("The WASM command can only be run from WebAssembly")
}
