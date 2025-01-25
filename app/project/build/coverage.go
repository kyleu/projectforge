package build

import (
	"context"

	"github.com/pkg/errors"

	"projectforge.dev/projectforge/app/lib/filesystem"
	"projectforge.dev/projectforge/app/util"
)

func Coverage(ctx context.Context, fs filesystem.FileLoader, scope string, logger util.Logger) (any, []string, error) {
	ex := &ExecHelper{}
	var logs []string

	if scope == "" {
		scope = "./app/..."
	}

	cmd := "go test -race -coverprofile ./tmp/coverage.out " + scope
	out, err := ex.Cmd(ctx, "coverage-run", cmd, fs.Root(), logger)
	if err != nil {
		return nil, ex.Logs, errors.Wrapf(err, "unable to run [%s]", cmd)
	}

	ret := util.StringSplitLines(out)
	return ret, logs, nil
}
