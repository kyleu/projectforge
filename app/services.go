// Package app
package app // import "projectforge.dev/projectforge/app"

import (
	"context"

	"projectforge.dev/projectforge/app/module"
	"projectforge.dev/projectforge/app/project"
	"projectforge.dev/projectforge/app/project/export"
	"projectforge.dev/projectforge/app/util"
)

type Services struct {
	CoreServices
	Modules  *module.Service
	Projects *project.Service
	Export   *export.Service
}

func NewServices(ctx context.Context, st *State, rootLogger util.Logger) (*Services, error) {
	ms, err := module.NewService(ctx, st.Files.Root(), rootLogger)
	if err != nil {
		return nil, err
	}
	ps, es := project.NewService(), export.NewService()
	core := initCoreServices(ctx, st, rootLogger)
	core.Socket.ReplaceHandlers(nil, nil)
	return &Services{CoreServices: core, Modules: ms, Projects: ps, Export: es}, nil
}

func (s *Services) Close(_ context.Context, _ util.Logger) error {
	return nil
}
