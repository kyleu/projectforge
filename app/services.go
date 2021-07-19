// Package app $PF_IGNORE$
package app

import (
	"github.com/kyleu/projectforge/app/module"
	"github.com/kyleu/projectforge/app/project"
)

type Services struct {
	Modules  *module.Service
	Projects *project.Service
}

func NewServices(st *State) (*Services, error) {
	return &Services{
		Modules:  module.NewService(st.Logger),
		Projects: project.NewService(st.Logger),
	}, nil
}
