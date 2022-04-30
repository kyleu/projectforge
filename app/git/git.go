package git

import (
	"context"

	"github.com/pkg/errors"
	"go.uber.org/zap"

	"projectforge.dev/projectforge/app/lib/telemetry"
)

var errNoRepo = errors.New("not a git repository")

func gitCmd(ctx context.Context, args string, path string, logger *zap.SugaredLogger) (string, error) {
	exit, out, err := telemetry.RunProcessSimple(ctx, "git "+args, path, logger)
	if err != nil {
		return "", errors.Wrap(err, "can't read git status for path ["+path+"]")
	}
	if exit == 128 {
		return "", errors.Wrapf(errNoRepo, "path [%s] is not a git repo", path)
	}
	if exit != 0 {
		return "", errors.Errorf("git status returned exit code [%d] for path [%s]: %s", exit, path, out)
	}
	return out, nil
}
