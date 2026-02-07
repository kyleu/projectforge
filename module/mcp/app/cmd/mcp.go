package cmd

import (
	"context"

	"github.com/muesli/coral"
	"github.com/pkg/errors"
	"go.uber.org/zap"

	"{{{ .Package }}}/app"
	"{{{ .Package }}}/app/lib/log"
	"{{{ .Package }}}/app/lib/mcpserver"
	"{{{ .Package }}}/app/util"
)

func mcpCmd() *coral.Command {
	f := func(*coral.Command, []string) error { return runMCP(rootCtx) }
	ret := newCmd("mcp", "Handles Model Context Protocol requests", f)
	return ret
}

func runMCP(ctx context.Context) error {
	// override logging
	l, err := log.InitDevLogging(log.GetLevel(zap.FatalLevel))
	if err != nil {
		return errors.Wrap(err, "error initializing logging")
	}
	logger := l.Sugar()
	util.RootLogger = logger

	if _, err = initIfNeeded(ctx); err != nil {
		return errors.Wrap(err, "error initializing application")
	}
	st, err := app.Bootstrap(ctx, _buildInfo, _flags.ConfigDir, _flags.Port, false, logger)
	if err != nil {
		return err
	}
	mcp, err := mcpserver.NewServer(ctx, st, logger)
	if err != nil {
		return err
	}
	if err := mcp.ServeCLI(ctx); err != nil {
		return err
	}
	return nil
}
