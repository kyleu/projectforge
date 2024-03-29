package git

import (
	"context"
	"fmt"
	"strings"

	"github.com/pkg/errors"
	"github.com/samber/lo"

	"projectforge.dev/projectforge/app/project"
	"projectforge.dev/projectforge/app/util"
)

func (s *Service) Push(ctx context.Context, prj *project.Project, logger util.Logger) (*Result, error) {
	_, err := gitFetch(ctx, prj.Path, false, logger)
	if err != nil {
		return nil, errors.Wrap(err, "unable to fetch for pull")
	}
	x, err := gitPush(ctx, prj.Path, logger)
	if err != nil {
		return nil, errors.Wrap(err, "unable to pull")
	}
	count := lo.CountBy(util.StringSplitLines(x), func(line string) bool {
		return strings.HasPrefix(line, "   ")
	})
	status := ok
	fetched := noUpdates
	if count > 0 {
		status = fmt.Sprintf("[%s] fetched", util.StringPlural(count, "update"))
		fetched = status
	}

	return NewResult(prj, status, util.ValueMap{"updates": fetched}), nil
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
