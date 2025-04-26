package git

import (
	"context"
	"fmt"
	"strings"

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
	_, err := GitCmd(ctx, "add .", path, logger)
	if err != nil {
		if isNotRepo(err) {
			return "", nil
		}
		return "", err
	}
	out, err := GitCmd(ctx, fmt.Sprintf(`commit -m %q`, message), path, logger)
	if err != nil {
		return "", err
	}
	return out, nil
}

func (s *Service) CommitCount(ctx context.Context, all bool, logger util.Logger) (*Result, error) {
	result, err := gitCommitCount(ctx, s.Path, all, logger)
	if err != nil {
		return nil, err
	}
	return NewResult(s.Key, util.OK, util.ValueMap{"count": result}), nil
}

func gitCommitCount(ctx context.Context, path string, all bool, logger util.Logger) (int, error) {
	cmd := "rev-list --count HEAD"
	if all {
		cmd = "rev-list --all --count"
	}
	out, err := GitCmd(ctx, cmd, path, logger)
	if err != nil {
		if isNotRepo(err) {
			return 0, nil
		}
		return 0, err
	}
	return util.ParseIntSimple(strings.TrimSpace(out)), nil
}
