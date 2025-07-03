package cmd

import (
	"context"
	"os"

	"projectforge.dev/projectforge/app"
	"projectforge.dev/projectforge/app/lib/log"
)

func Entrypoint(ctx context.Context, bi *app.BuildInfo) {
	logger, err := Run(ctx, bi)
	if err != nil {
		const msg = "exiting due to error"
		if logger == nil {
			println(log.Red.Add(err.Error())) //nolint:forbidigo
			println(log.Red.Add(msg))         //nolint:forbidigo
		} else {
			logger.Error(err)
			logger.Error(msg)
		}
		os.Exit(1)
	}
}
