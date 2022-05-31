// Package app $PF_IGNORE$
package app // import "projectforge.dev/projectforge/app"

import (
	"context"

	"projectforge.dev/projectforge/app/export"
	"projectforge.dev/projectforge/app/git"
	"projectforge.dev/projectforge/app/module"
	"projectforge.dev/projectforge/app/project"
	"projectforge.dev/projectforge/app/util"
)

type Services struct {
	Modules  *module.Service
	Projects *project.Service
	Export   *export.Service
	Git      *git.Service
}

func NewServices(ctx context.Context, st *State, rootLogger util.Logger) (*Services, error) {
	return &Services{
		Modules:  module.NewService(ctx, st.Files, rootLogger),
		Projects: project.NewService(),
		Export:   export.NewService(),
		Git:      git.NewService(),
	}, nil
}

func (s *Services) Close(_ context.Context, _ util.Logger) error {
	return nil
}
