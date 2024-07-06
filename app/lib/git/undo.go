package git

import (
	"context"
	"fmt"

	"projectforge.dev/projectforge/app/util"
)

func (s *Service) UndoCommit(ctx context.Context, prj string, path string, logger util.Logger) (*Result, error) {
	result, err := gitResetSoft(ctx, path, logger)
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
