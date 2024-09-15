package git

import (
	"context"
	"fmt"
	"strings"

	"github.com/pkg/errors"
	"github.com/samber/lo"

	"{{{ .Package }}}/app/util"
)

func (s *Service) Reset(ctx context.Context, logger util.Logger) (*Result, error) {
	x, err := gitReset(ctx, s.Path, logger)
	if err != nil {
		return nil, errors.Wrap(err, "unable to reset git repo")
	}
	count := lo.CountBy(util.StringSplitLines(x), func(line string) bool {
		return strings.HasPrefix(line, "   ")
	})
	status := util.OK
	fetched := "no changes"
	if count > 0 {
		status = fmt.Sprintf("[%s] reset", util.StringPlural(count, "file"))
		fetched = status
	}

	return NewResult(s.Key, status, util.ValueMap{"updates": fetched}), nil
}

func gitReset(ctx context.Context, path string, logger util.Logger) (string, error) {
	out, err := gitCmd(ctx, "reset --hard", path, logger)
	if err != nil {
		if isNoRepo(err) {
			return "", nil
		}
		return "", err
	}
	return out, nil
}
