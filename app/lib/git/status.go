package git

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/pkg/errors"

	"projectforge.dev/projectforge/app/lib/telemetry"
	"projectforge.dev/projectforge/app/util"
)

func (s *Service) Status(ctx context.Context, logger util.Logger) (*Result, error) {
	_, span, _ := telemetry.StartSpan(ctx, "git.status:"+s.Key, logger)
	defer span.Complete()

	data := make(util.ValueMap, 16)

	commitsAhead, commitsBehind, dirty, err := gitStatus(ctx, s.Path, logger)
	if err != nil {
		return nil, errors.Wrap(err, "unable to find git status")
	}
	if commitsAhead != 0 {
		data["commitsAhead"] = commitsAhead
	}
	if commitsBehind != 0 {
		data["commitsBehind"] = commitsBehind
	}

	data["branch"] = gitBranch(ctx, s.Path, logger)
	if len(dirty) > 0 {
		data["dirty"] = dirty
	}
	status := util.OK
	if len(dirty) > 0 {
		status = fmt.Sprintf("[%d] changes", len(dirty))
	}
	return NewResult(s.Key, status, data), nil
}

func gitStatus(ctx context.Context, path string, logger util.Logger) (int, int, []string, error) {
	out, err := GitCmd(ctx, "status --porcelain -b", path, logger)
	if err != nil {
		if isNotRepo(err) {
			return 0, 0, nil, nil
		}
		return 0, 0, nil, err
	}

	lines := util.StringSplitAndTrim(out, util.StringDetectLinebreak(out))

	var commitsAhead, commitsBehind int
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
				ca, err := strconv.ParseInt(l, 10, 32)
				if err != nil {
					return 0, 0, nil, errors.Wrap(err, "invalid [ahead] block")
				}
				commitsAhead = int(ca)
			}
			if i := strings.Index(line, "behind "); i > -1 {
				l := line[i+7:]
				end := strings.Index(l, "]")
				l = strings.TrimSpace(l[:end])
				cb, err := strconv.ParseInt(l, 10, 32)
				if err != nil {
					return 0, 0, nil, errors.Wrap(err, "invalid [behind] block")
				}
				commitsBehind = int(cb)
			}
			continue
		}
		if i := strings.Index(line, " "); i > -1 {
			line = line[i+1:]
		}
		dirty = append(dirty, strings.TrimSpace(line))
	}

	return commitsAhead, commitsBehind, util.ArraySorted(dirty), nil
}
