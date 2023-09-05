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

func (s *Service) Fetch(ctx context.Context, prj *project.Project, logger util.Logger) (*Result, error) {
	x, err := gitFetch(ctx, prj.Path, true, logger)
	if err != nil {
		return nil, errors.Wrap(err, "unable to fetch")
	}
	count := lo.CountBy(util.StringSplitLines(x), func(line string) bool {
		return strings.HasPrefix(line, "   ")
	})
	status := ok
	fetched := noUpdates
	if count > 0 {
		status = fmt.Sprintf("[%d] %s fetched", count, util.StringPlural(count, "update"))
		fetched = status
	}

	return NewResult(prj, status, util.ValueMap{"updates": fetched}), nil
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
