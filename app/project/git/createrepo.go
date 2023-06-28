package git

import (
	"context"

	"github.com/pkg/errors"

	"projectforge.dev/projectforge/app/project"
	"projectforge.dev/projectforge/app/util"
)

func (s *Service) CreateRepo(ctx context.Context, prj *project.Project, logger util.Logger) (*Result, error) {
	err := gitCreateRepo(ctx, prj.Path, logger)
	if err != nil {
		return nil, errors.Wrap(err, "unable to create repo")
	}
	return NewResult(prj, ok, util.ValueMap{"repo": "created"}), nil
}

func gitCreateRepo(ctx context.Context, path string, logger util.Logger) error {
	_, err := gitCmd(ctx, "git init", path, logger)
	if err != nil {
		if isNoRepo(err) {
			return nil
		}
		return err
	}
	return nil
}
