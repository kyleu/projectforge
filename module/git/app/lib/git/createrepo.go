package git

import (
	"context"

	"github.com/pkg/errors"

	"{{{ .Package }}}/app/util"
)

func (s *Service) CreateRepo(ctx context.Context, logger util.Logger) (*Result, error) {
	err := gitCreateRepo(ctx, s.Path, logger)
	if err != nil {
		return nil, errors.Wrap(err, "unable to create repo")
	}
	return NewResult(s.Key, util.OK, util.ValueMap{"repo": "created"}), nil
}

func gitCreateRepo(ctx context.Context, path string, logger util.Logger) error {
	_, err := gitCmd(ctx, "git init", path, logger)
	if err != nil {
		if isNotRepo(err) {
			return nil
		}
		return err
	}
	return nil
}
