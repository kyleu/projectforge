package cmd

import (
	"context"

	"github.com/muesli/coral"
	"github.com/pkg/errors"

	"projectforge.dev/projectforge/app/controller/tui"
	"projectforge.dev/projectforge/app/util"
)

func tuiCmd() *coral.Command {
	f := func(*coral.Command, []string) error { return runTUI(rootCtx) }
	return newCmd("tui", "Interactive CLI application", f)
}

func runTUI(ctx context.Context) error {
	logger, err := initIfNeeded(ctx)
	if err != nil {
		return errors.Wrap(err, "error initializing application")
	}
	st, r, logger, err := loadServer(ctx, _flags, logger)
	if err != nil {
		return err
	}
	ss := func() error {
		_, err := listenAndServe(ctx, _flags.Address, _flags.Port, r)
		return err
	}

	ret, err := tui.RunTUI(ctx, st, ss, logger)
	if err != nil {
		return err
	}
	logger.Debugf("concluded TUI session with result [%s]", util.ToJSON(ret))
	return nil
}
