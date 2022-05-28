package git

import (
	"context"
	"fmt"
	"strings"

	"github.com/pkg/errors"
	"golang.org/x/exp/slices"
	"projectforge.dev/projectforge/app/lib/telemetry"
	"projectforge.dev/projectforge/app/project"
	"projectforge.dev/projectforge/app/util"
)

func (s *Service) Status(ctx context.Context, prj *project.Project, logger util.Logger) (*Result, error) {
	_, span, _ := telemetry.StartSpan(ctx, "git.status:"+prj.Key, logger)
	defer span.Complete()

	dirty, err := gitStatus(ctx, prj.Path, logger)
	if err != nil {
		return nil, errors.Wrap(err, "unable to find git status")
	}
	branch := gitBranch(ctx, prj.Path, logger)
	data := util.ValueMap{"branch": branch}
	if len(dirty) > 0 {
		data["dirty"] = dirty
	}
	status := ok
	if len(dirty) > 0 {
		status = fmt.Sprintf("[%d] changes", len(dirty))
	}
	return NewResult(prj, status, data), nil
}

func gitStatus(ctx context.Context, path string, logger util.Logger) ([]string, error) {
	out, err := gitCmd(ctx, "status --porcelain", path, logger)
	if err != nil {
		if errors.Is(err, errNoRepo) {
			return nil, nil
		}
		return nil, err
	}

	lines := util.StringSplitAndTrim(out, "\n")

	dirty := make([]string, 0, len(lines))
	for _, line := range lines {
		if i := strings.Index(line, " "); i > -1 {
			line = line[i+1:]
		}
		dirty = append(dirty, strings.TrimSpace(line))
	}
	slices.Sort(dirty)

	return dirty, nil
}
