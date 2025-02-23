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
	mcp, err := mcpserver.NewServer(ctx, _version)
	if err != nil {
		return err
	}
	// $PF_SECTION_START(tools)$
	// $PF_SECTION_END(tools)$
	if err := mcp.Serve(ctx); err != nil {
		return err
	}
	return nil
}
