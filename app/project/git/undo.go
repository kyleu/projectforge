package git

import (
	"context"
	"fmt"

	"projectforge.dev/projectforge/app/project"
	"projectforge.dev/projectforge/app/util"
)

func (s *Service) UndoCommit(ctx context.Context, prj *project.Project, logger util.Logger) (*Result, error) {
	result, err := gitResetSoft(ctx, prj.Path, logger)
	if err != nil {
		return nil, err
	}

	return NewResult(prj, ok, util.ValueMap{"reset": result}), nil
}

func gitResetSoft(ctx context.Context, path string, logger util.Logger) (string, error) {
	currBranch := gitBranch(ctx, path, logger)
	_, err := gitCmd(ctx, fmt.Sprintf("reset --soft %s~1", currBranch), path, logger)
	if err != nil {
		if isNoRepo(err) {
			return "", nil
		}
		return "", err
	}
	return "OK", nil
}
