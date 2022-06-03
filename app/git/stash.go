package git

import (
	"context"

	"github.com/pkg/errors"

	"projectforge.dev/projectforge/app/project"
	"projectforge.dev/projectforge/app/util"
)

func (s *Service) gitStash(ctx context.Context, prj *project.Project, logger util.Logger) (string, error) {
	out, err := gitCmd(ctx, "stash", prj.Path, logger)
	if err != nil {
		if isNoRepo(err) {
			return "", nil
		}
		return "", errors.Wrap(err, "unable to apply stash")
	}
	return out, nil
}

func (s *Service) gitStashPop(ctx context.Context, prj *project.Project, logger util.Logger) (string, error) {
	out, err := gitCmd(ctx, "stash pop", prj.Path, logger)
	if err != nil {
		if isNoRepo(err) {
			return "", nil
		}
		return "", errors.Wrap(err, "unable to pop stash")
	}
	return out, nil
}
