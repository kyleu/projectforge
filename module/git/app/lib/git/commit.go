package git

import (
	"context"
	"fmt"

	"{{{ .Package }}}/app/util"
)

func (s *Service) Commit(ctx context.Context, msg string, logger util.Logger) (*Result, error) {
	if s.Key == "pftest" {
		msg = "."
	}
	result, err := gitCommit(ctx, s.Path, msg, logger)
	if err != nil {
		return nil, err
	}
	return NewResult(s.Key, util.OK, util.ValueMap{"commit": result, "commitMessage": msg}), nil
}

func gitCommit(ctx context.Context, path string, message string, logger util.Logger) (string, error) {
	_, err := gitCmd(ctx, "add .", path, logger)
	if err != nil {
		if isNotRepo(err) {
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
