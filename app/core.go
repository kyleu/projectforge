package app

import (
	"context"

	"projectforge.dev/projectforge/app/lib/exec"
	"projectforge.dev/projectforge/app/lib/help"
	"projectforge.dev/projectforge/app/lib/task"
	"projectforge.dev/projectforge/app/lib/websocket"
	"projectforge.dev/projectforge/app/util"
)

type CoreServices struct {
	Exec   *exec.Service
	Socket *websocket.Service
	Task   *task.Service
	Help   *help.Service
}

func initCoreServices(ctx context.Context, st *State, logger util.Logger) CoreServices {
	return CoreServices{
		Exec:   exec.NewService(),
		Socket: websocket.NewService(nil, nil),
		Task:   task.NewService(st.Files, "task_history"),
		Help:   help.NewService(logger),
	}
}
