package git

import (
	"context"
	"fmt"

	"projectforge.dev/projectforge/app/project"
	"projectforge.dev/projectforge/app/util"
)

func (s *Service) Commit(ctx context.Context, prj *project.Project, msg string, logger util.Logger) (*Result, error) {
	result, err := gitCommit(ctx, prj.Path, msg, logger)
	if err != nil {
		return nil, err
	}

	return NewResult(prj, ok, util.ValueMap{"commit": result, "commitMessage": msg}), nil
}

func gitCommit(ctx context.Context, path string, message string, logger util.Logger) (string, error) {
	_, err := gitCmd(ctx, "add .", path, logger)
	if err != nil {
		if isNoRepo(err) {
			return "", nil
		}
		return "", err
	}
	out, err := gitCmd(ctx, fmt.Sprintf(`commit -m %q`, message), path, logger)
	if err != nil {
		return "", err
	}
	return out, nil
}
