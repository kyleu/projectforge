package mcpserver

import (
	"context"

	"projectforge.dev/projectforge/app"
	"projectforge.dev/projectforge/app/util"
)

var defaultServer *Server
var defaultTools Tools

func AddDefaultTools(ts ...*Tool) {
	defaultTools = append(defaultTools, ts...)
}

func ClearDefaultTools() {
	defaultTools = nil
}

func GetDefaultServer(ctx context.Context, as *app.State, logger util.Logger) (*Server, error) {
	if defaultServer != nil {
		return defaultServer, nil
	}
	ret, err := NewServer(ctx, as, logger)
	if err != nil {
		return nil, err
	}
	defaultServer = ret
	return defaultServer, nil
}
