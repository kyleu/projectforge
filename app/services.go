// Package app $PF_IGNORE$
package app // import "projectforge.dev/app"

import (
	"context"

	"projectforge.dev/app/git"
	"projectforge.dev/app/module"
	"projectforge.dev/app/project"
	"projectforge.dev/app/export"
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
