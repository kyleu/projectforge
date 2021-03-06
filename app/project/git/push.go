package git

import (
	"context"
	"fmt"
	"strings"

	"github.com/pkg/errors"

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
	count := 0
	lines := strings.Split(x, "\n")
	for _, line := range lines {
		if strings.HasPrefix(line, "   ") {
			count++
		}
	}
	status := ok
	fetched := noUpdates
	if count > 0 {
		status = fmt.Sprintf("[%d] %s fetched", count, util.StringPluralMaybe("update", count))
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
