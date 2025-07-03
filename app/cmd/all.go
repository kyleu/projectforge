package cmd

import (
	"context"
	"fmt"

	"github.com/muesli/coral"
	"github.com/pkg/errors"

	"projectforge.dev/projectforge/app/util"
)

const keyAll = "all"

func allCmd() *coral.Command {
	short := fmt.Sprintf("Starts the main http server on port %d and the marketing site on port %d", util.AppPort, util.AppPort+1)
	f := func(*coral.Command, []string) error { return allF(rootCtx) }
	ret := &coral.Command{Use: keyAll, Short: short, RunE: f}
	return ret
}

func allF(ctx context.Context) error {
	logger, err := initIfNeeded(ctx)
	if err != nil {
		return errors.Wrap(err, "error initializing application")
	}

	go func() {
		if err := startSite(ctx, _flags.Clone(_flags.Port+1)); err != nil {
			logger.Errorf("unable to start marketing site: %+v", err)
		}
	}()
	return startServer(ctx, _flags)
}
