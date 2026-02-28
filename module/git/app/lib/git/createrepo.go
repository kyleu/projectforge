package git

import (
	"context"
	"fmt"

	"github.com/pkg/errors"

	"{{{ .Package }}}/app/util"
)

func (s *Service) CreateRepo(ctx context.Context, logger util.Logger) (*Result, error) {
	err := gitCreateRepo(ctx, s.Path, logger)
	if err != nil {
		return nil, errors.Wrap(err, "unable to create repo")
	}
	return NewResult(s.Key, util.KeyOK, util.ValueMap{"repo": "created"}), nil
}

func gitCreateRepo(ctx context.Context, path string, logger util.Logger) error {
	_, err := GitCmd(ctx, "git init", path, logger)
	if err != nil {
		if isNotRepo(err) {
			return nil
		}
		return err
	}
	return nil
}

func CloneRepo(ctx context.Context, path string, url string, logger util.Logger) error {
	_, err := GitCmd(ctx, fmt.Sprintf("clone %q", url), path, logger)
	if err != nil {
		if isNotRepo(err) {
			return nil
		}
		return err
	}
	return nil
}

func CloneRepoGH(ctx context.Context, path string, url string, logger util.Logger) error {
	_, err := GHCmd(ctx, fmt.Sprintf("repo clone %q", url), path, logger)
	if err != nil {
		if isNotRepo(err) {
			return nil
		}
		return err
	}
	return nil
}
