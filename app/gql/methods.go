package gql

import (
	"context"
	"github.com/samber/lo"
	"projectforge.dev/projectforge/app/module"
	"projectforge.dev/projectforge/app/project"
)

func (s *Schema) Hello(ctx context.Context) (string, error) {
	return "Howdy!", nil
}

func (s *Schema) Modules(ctx context.Context) ([]*Module, error) {
	return lo.Map(s.st.Services.Modules.ModulesSorted(), func(x *module.Module, _ int) *Module {
		return FromModule(x, s.log)
	}), nil
}

func (s *Schema) Projects(ctx context.Context) ([]*Project, error) {
	return lo.Map(s.st.Services.Projects.Projects(), func(x *project.Project, _ int) *Project {
		return FromProject(x, s.log)
	}), nil
}

func (s *Schema) Poke(ctx context.Context) (string, error) {
	return "OK!", nil
}
