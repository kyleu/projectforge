// Package app $PF_IGNORE$
package app // import "projectforge.dev/projectforge/app"

import (
	"context"
	"encoding/json"

	"projectforge.dev/projectforge/app/lib/exec"
	"projectforge.dev/projectforge/app/lib/websocket"
	"projectforge.dev/projectforge/app/module"
	"projectforge.dev/projectforge/app/project"
	"projectforge.dev/projectforge/app/project/export"
	"projectforge.dev/projectforge/app/project/git"
	"projectforge.dev/projectforge/app/util"
)

type Services struct {
	Modules  *module.Service
	Projects *project.Service
	Export   *export.Service
	Git      *git.Service
	Exec     *exec.Service
	Socket   *websocket.Service
}

func NewServices(ctx context.Context, st *State, rootLogger util.Logger) (*Services, error) {
	return &Services{
		Modules:  module.NewService(ctx, st.Files, rootLogger),
		Projects: project.NewService(),
		Export:   export.NewService(),
		Git:      git.NewService(),
		Exec:     exec.NewService(),
		Socket:   websocket.NewService(rootLogger, nil, socketHandler, nil, nil),
	}, nil
}

func (s *Services) Close(_ context.Context, _ util.Logger) error {
	return nil
}

func socketHandler(s *websocket.Service, c *websocket.Connection, svc string, cmd string, param json.RawMessage) error {
	switch cmd {
	case "connect":
		_, err := s.Join(c.ID, "tap")
		if err != nil {
			return err
		}
	default:
		s.Logger.Error("unhandled command [" + cmd + "]")
	}
	return nil
}
