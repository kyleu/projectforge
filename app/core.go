// Package app - Content managed by Project Forge, see [projectforge.md] for details.
package app

import (
	"context"

	"projectforge.dev/projectforge/app/lib/exec"
	"projectforge.dev/projectforge/app/lib/help"
	"projectforge.dev/projectforge/app/lib/websocket"
	"projectforge.dev/projectforge/app/util"
)

type CoreServices struct {
	Exec   *exec.Service
	Socket *websocket.Service
	Help   *help.Service
}

func initCoreServices(ctx context.Context, st *State, logger util.Logger) CoreServices {
	return CoreServices{
		Exec:   exec.NewService(),
		Socket: websocket.NewService(nil, nil, nil),
		Help:   help.NewService(logger),
	}
}
