// Package app
package app // import "projectforge.dev/projectforge/app"

import (
	"context"
	"encoding/json"

	"projectforge.dev/projectforge/app/lib/websocket"
	"projectforge.dev/projectforge/app/module"
	"projectforge.dev/projectforge/app/project"
	"projectforge.dev/projectforge/app/project/export"
	"projectforge.dev/projectforge/app/project/git"
	"projectforge.dev/projectforge/app/util"
)

type Services struct {
	CoreServices
	Modules  *module.Service
	Projects *project.Service
	Export   *export.Service
	Git      *git.Service
}

func NewServices(ctx context.Context, st *State, rootLogger util.Logger) (*Services, error) {
	ms, err := module.NewService(ctx, st.Files.Root(), rootLogger)
	if err != nil {
		return nil, err
	}
	ps, es, gs := project.NewService(), export.NewService(), git.NewService()
	core := initCoreServices(ctx, st, rootLogger)
	core.Socket.ReplaceHandlers(nil, socketHandler, nil)
	return &Services{CoreServices: core, Modules: ms, Projects: ps, Export: es, Git: gs}, nil
}

func (s *Services) Close(_ context.Context, _ util.Logger) error {
	return nil
}

func socketHandler(_ context.Context, s *websocket.Service, c *websocket.Connection, _ string, cmd string, _ json.RawMessage, logger util.Logger) error {
	switch cmd {
	case "connect":
		_, err := s.Join(c.ID, "tap", logger)
		if err != nil {
			return err
		}
	default:
		logger.Error("unhandled command [" + cmd + "]")
	}
	return nil
}
