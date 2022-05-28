package git

import (
	"context"

	"projectforge.dev/projectforge/app/project"
	"projectforge.dev/projectforge/app/util"
)

func (s *Service) CreateRepo(ctx context.Context, prj *project.Project, logger util.Logger) (*Result, error) {
	return NewResult(prj, "TODO", util.ValueMap{"TODO": "Create Repo"}), nil
}
