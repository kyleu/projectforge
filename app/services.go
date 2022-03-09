// Package app $PF_IGNORE$
package app // import "projectforge.dev/projectforge/app"

import (
	"context"

	"projectforge.dev/projectforge/app/export"
	"projectforge.dev/projectforge/app/git"
	"projectforge.dev/projectforge/app/module"
	"projectforge.dev/projectforge/app/project"
)

type Services struct {
	Modules  *module.Service
	Projects *project.Service
	Export   *export.Service
	Git      *git.Service
}

func NewServices(ctx context.Context, st *State) (*Services, error) {
	return &Services{
		Modules:  module.NewService(ctx, st.Files, st.Logger),
		Projects: project.NewService(st.Logger),
		Export:   export.NewService(st.Logger),
		Git:      git.NewService(st.Logger),
	}, nil
}
