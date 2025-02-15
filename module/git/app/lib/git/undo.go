package git

import (
	"context"
	"fmt"

	"{{{ .Package }}}/app/util"
)

func (s *Service) UndoCommit(ctx context.Context, logger util.Logger) (*Result, error) {
	result, err := gitResetSoft(ctx, s.Path, logger)
	if err != nil {
		return nil, err
	}

	return NewResult(s.Key, util.OK, util.ValueMap{"reset": result}), nil
}

func gitResetSoft(ctx context.Context, path string, logger util.Logger) (string, error) {
	currBranch := gitBranch(ctx, path, logger)
	_, err := gitCmd(ctx, fmt.Sprintf("reset --soft %s~1", currBranch), path, logger)
	if err != nil {
		if isNotRepo(err) {
			return "", nil
		}
		return "", err
	}
	return util.OK, nil
}
