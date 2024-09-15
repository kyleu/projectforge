package git

import (
	"context"
	"fmt"
	"strings"

	"github.com/pkg/errors"
	"github.com/samber/lo"

	"projectforge.dev/projectforge/app/util"
)

func (s *Service) Push(ctx context.Context, logger util.Logger) (*Result, error) {
	_, err := gitFetch(ctx, s.Path, false, logger)
	if err != nil {
		return nil, errors.Wrap(err, "unable to fetch for pull")
	}
	x, err := gitPush(ctx, s.Path, logger)
	if err != nil {
		return nil, errors.Wrap(err, "unable to pull")
	}
	count := lo.CountBy(util.StringSplitLines(x), func(line string) bool {
		return strings.HasPrefix(line, "   ")
	})
	status := util.OK
	fetched := noUpdates
	if count > 0 {
		status = fmt.Sprintf("[%s] pushed", util.StringPlural(count, "commit"))
		fetched = status
	}

	return NewResult(s.Key, status, util.ValueMap{"updates": fetched}), nil
}

func gitPush(ctx context.Context, path string, logger util.Logger) (string, error) {
	out, err := gitCmd(ctx, "push", path, logger)
	if err != nil {
		if isNoRepo(err) {
			return "", nil
		}
		return "", err
	}
	return out, nil
}
