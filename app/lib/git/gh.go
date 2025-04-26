package git

import (
	"context"

	"github.com/pkg/errors"

	"projectforge.dev/projectforge/app/lib/telemetry"
	"projectforge.dev/projectforge/app/util"
)

func GHCmd(ctx context.Context, args string, path string, logger util.Logger) (string, error) {
	exit, out, err := telemetry.RunProcessSimple(ctx, "gh "+args, path, logger)
	if err != nil {
		return "", errors.Wrapf(err, "can't run [gh %s] for path [%s]", args, path)
	}
	if exit != 0 {
		return out, errors.Errorf("gh cmd [gh %s] returned exit code [%d] for path [%s]: %s", args, exit, path, out)
	}
	return out, nil
}
