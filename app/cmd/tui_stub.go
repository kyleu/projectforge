//go:build js
// +build js

package cmd

import (
	"context"
	"fmt"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"

	"projectforge.dev/projectforge/app/util"
)

const keyTUI = "tui"

func tuiCmd() *cobra.Command {
	short := fmt.Sprintf("Starts the terminal UI (and the http server on port %d)", util.AppPort)
	f := func(*cobra.Command, []string) error { return runTUI(rootCtx, _flags) }
	ret := &cobra.Command{Use: keyTUI, Short: short, RunE: f}
	return ret
}

func runTUI(ctx context.Context, flags *Flags) error {
	return errors.New("The [tui] command can't be run from WebAssembly")
}
