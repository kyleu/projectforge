package cmd

import (
	"context"

	"github.com/muesli/coral"

	"{{{ .Package }}}/app/lib/mcpserver"
)

func mcpCmd() *coral.Command {
	f := func(*coral.Command, []string) error { return runMCP(context.Background()) }
	ret := &coral.Command{Use: "mcp", Short: "Handles Model Context Protocol requests", RunE: f}
	return ret
}

func runMCP(ctx context.Context) error {
	if err := initIfNeeded(); err != nil {
		return err
	}
	mcpserver.InitMCP(_buildInfo, _flags.Debug)
	mcp, err := mcpserver.NewServer(ctx, _buildInfo.Version)
	if err != nil {
		return err
	}
	if err := mcp.Serve(ctx); err != nil {
		return err
	}
	return nil
}
