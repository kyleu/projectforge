package git

import (
	"context"
	"fmt"
	"strconv"
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

	data := make(util.ValueMap, 16)

	commitsAhead, commitsBehind, dirty, err := gitStatus(ctx, prj.Path, logger)
	if err != nil {
		return nil, errors.Wrap(err, "unable to find git status")
	}
	if commitsAhead != 0 {
		data["commitsAhead"] = commitsAhead
	}
	if commitsBehind != 0 {
		data["commitsBehind"] = commitsBehind
	}

	data["branch"] = gitBranch(ctx, prj.Path, logger)
	if len(dirty) > 0 {
		data["dirty"] = dirty
	}
	status := ok
	if len(dirty) > 0 {
		status = fmt.Sprintf("[%d] changes", len(dirty))
	}
	return NewResult(prj, status, data), nil
}

func gitStatus(ctx context.Context, path string, logger util.Logger) (int, int, []string, error) {
	out, err := gitCmd(ctx, "status --porcelain -b", path, logger)
	if err != nil {
		if isNoRepo(err) {
			return 0, 0, nil, nil
		}
		return 0, 0, nil, err
	}

	lines := util.StringSplitAndTrim(out, "\n")

	commitsAhead, commitsBehind := 0, 0
	dirty := make([]string, 0, len(lines))
	for _, line := range lines {
		if strings.HasPrefix(line, "##") {
			if i := strings.Index(line, "ahead "); i > -1 {
				l := line[i+6:]
				end := strings.Index(l, ",")
				if end == -1 {
					end = strings.Index(l, "]")
				}
				l = strings.TrimSpace(l[:end])
				commitsAhead, err = strconv.Atoi(l)
				if err != nil {
					return 0, 0, nil, errors.Wrap(err, "invalid [ahead] block")
				}
			}
			if i := strings.Index(line, "behind "); i > -1 {
				l := line[i+7:]
				end := strings.Index(l, "]")
				l = strings.TrimSpace(l[:end])
				commitsBehind, err = strconv.Atoi(l)
				if err != nil {
					return 0, 0, nil, errors.Wrap(err, "invalid [behind] block")
				}
			}
			continue
		}
		if i := strings.Index(line, " "); i > -1 {
			line = line[i+1:]
		}
		dirty = append(dirty, strings.TrimSpace(line))
	}
	slices.Sort(dirty)

	return commitsAhead, commitsBehind, dirty, nil
}
