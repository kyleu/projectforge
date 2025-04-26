package git

import (
	"context"
	"strings"

	"github.com/pkg/errors"

	"projectforge.dev/projectforge/app/lib/telemetry"
	"projectforge.dev/projectforge/app/util"
)

func isNotRepo(err error) bool {
	if err == nil {
		return false
	}
	return strings.Contains(err.Error(), "not a git repository")
}

func GitCmd(ctx context.Context, args string, path string, logger util.Logger) (string, error) {
	exit, out, err := telemetry.RunProcessSimple(ctx, "git "+args, path, logger)
	if err != nil {
		return "", errors.Wrapf(err, "can't run [git %s] for path [%s]", args, path)
	}
	if exit != 0 {
		return out, errors.Errorf("git cmd [git %s] returned exit code [%d] for path [%s]: %s", args, exit, path, out)
	}
	return out, nil
}
