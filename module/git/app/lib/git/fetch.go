package git

import (
	"context"
	"fmt"
	"strings"

	"github.com/pkg/errors"
	"github.com/samber/lo"

	"{{{ .Package }}}/app/util"
)

func (s *Service) Fetch(ctx context.Context, logger util.Logger) (*Result, error) {
	x, err := gitFetch(ctx, s.Path, true, logger)
	if err != nil {
		return nil, errors.Wrap(err, "unable to fetch")
	}
	count := lo.CountBy(util.StringSplitLines(x), func(line string) bool {
		return strings.HasPrefix(line, "   ")
	})
	status := util.OK
	fetched := noUpdates
	if count > 0 {
		status = fmt.Sprintf("[%s] fetched", util.StringPlural(count, "update"))
		fetched = status
	}

	return NewResult(s.Key, status, util.ValueMap{"updates": fetched}), nil
}

func gitFetch(ctx context.Context, path string, dryRun bool, logger util.Logger) (string, error) {
	cmd := "fetch"
	if dryRun {
		cmd += " --dry-run"
	}
	out, err := gitCmd(ctx, cmd, path, logger)
	if err != nil {
		if isNoRepo(err) {
			return "", nil
		}
		return "", err
	}
	return out, nil
}
