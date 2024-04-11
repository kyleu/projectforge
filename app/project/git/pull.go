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

func (s *Service) Pull(ctx context.Context, prj *project.Project, logger util.Logger) (*Result, error) {
	_, err := gitFetch(ctx, prj.Path, false, logger)
	if err != nil {
		return nil, errors.Wrap(err, "unable to fetch for pull")
	}
	x, err := gitPull(ctx, prj.Path, logger)
	if err != nil {
		return nil, errors.Wrap(err, "unable to pull")
	}
	count := lo.CountBy(util.StringSplitLines(x), func(line string) bool {
		return strings.HasPrefix(line, "   ")
	})
	status := ok
	fetched := noUpdates
	if count > 0 {
		status = fmt.Sprintf("[%s] pulled", util.StringPlural(count, "commit"))
		fetched = status
	}

	return NewResult(prj, status, util.ValueMap{"updates": fetched}), nil
}

func gitPull(ctx context.Context, path string, logger util.Logger) (string, error) {
	out, err := gitCmd(ctx, "pull -q", path, logger)
	if err != nil {
		if isNoRepo(err) {
			return "", nil
		}
		if strings.Contains(out, "divergent branches") {
			return "", errors.Errorf("there are conflicting changes to resolve before you can pull from [%s]", path)
		}
		return "", err
	}
	return out, nil
}
