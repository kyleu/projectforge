package tui

import (
	"context"

	"projectforge.dev/projectforge/app"
	"projectforge.dev/projectforge/app/util"
)

func RunTUI(ctx context.Context, st *app.State, startServerFn func() error, logger util.Logger) (any, error) {
	logger.Infof("TODO: run TUI")
	return "TODO", nil
}
