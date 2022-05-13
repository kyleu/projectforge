package git

import (
	"context"
	"fmt"
	"strings"

	"github.com/pkg/errors"
	"golang.org/x/exp/slices"
	"projectforge.dev/projectforge/app/util"
)

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

func gitBranch(ctx context.Context, path string, logger util.Logger) string {
	out, err := gitCmd(ctx, "branch --show-current", path, logger)
	if err != nil {
		if errors.Is(err, errNoRepo) {
			return "norepo"
		}
		return "error: " + err.Error()
	}
	return strings.TrimSpace(out)
}

func gitFetch(ctx context.Context, path string, dryRun bool, logger util.Logger) (string, error) {
	cmd := "fetch"
	if dryRun {
		cmd += " --dry-run"
	}
	out, err := gitCmd(ctx, cmd, path, logger)
	if err != nil {
		if errors.Is(err, errNoRepo) {
			return "", nil
		}
		return "", err
	}
	return out, nil
}

func gitCommit(ctx context.Context, path string, message string, logger util.Logger) (string, error) {
	_, err := gitCmd(ctx, "add .", path, logger)
	if err != nil {
		if errors.Is(err, errNoRepo) {
			return "", nil
		}
		return "", err
	}
	out, err := gitCmd(ctx, fmt.Sprintf(`commit -m %q`, message), path, logger)
	if err != nil {
		return "", err
	}
	return out, nil
}
