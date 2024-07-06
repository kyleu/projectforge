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

func (s *Service) Outdated(ctx context.Context, prj string, path string, logger util.Logger) (*Result, error) {
	_, span, _ := telemetry.StartSpan(ctx, "git.outdated:"+prj, logger)
	defer span.Complete()

	data := make(util.ValueMap, 16)

	tag, commitsAhead, err := gitOutdated(ctx, path, logger)
	if err != nil {
		return nil, errors.Wrap(err, "unable to find git outdated commits")
	}
	data["tag"] = tag
	if commitsAhead != 0 {
		data["commitsAhead"] = commitsAhead
	}
	status := tag
	if commitsAhead > 0 {
		status += fmt.Sprintf(" - %d commits ahead", commitsAhead)
	}
	return NewResult(prj, status, data), nil
}

func gitOutdated(ctx context.Context, path string, logger util.Logger) (string, int, error) {
	out, err := gitCmd(ctx, "rev-parse HEAD", path, logger)
	if err != nil {
		if isNoRepo(err) {
			return "", 0, nil
		}
		return "", 0, err
	}
	currentCommitHash := strings.TrimSpace(out)

	out, err = gitCmd(ctx, "describe --abbrev=0 --tags "+currentCommitHash, path, logger)
	if err != nil {
		return "", 0, nil //nolint:nilerr
	}
	latestTag := strings.TrimSpace(out)
	out, err = gitCmd(ctx, "rev-list --count "+fmt.Sprintf("%s..HEAD", latestTag), path, logger)
	if err != nil {
		return "", 0, err
	}
	numCommits, err := strconv.ParseInt(strings.TrimSpace(out), 10, 32)
	if err != nil {
		return "", 0, err
	}
	return latestTag, int(numCommits), nil
}
