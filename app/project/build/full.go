package build

import (
	"context"
	"fmt"
	"path/filepath"

	"github.com/pkg/errors"

	"projectforge.dev/projectforge/app/lib/telemetry"
	"projectforge.dev/projectforge/app/project"
	"projectforge.dev/projectforge/app/util"
)

func Full(ctx context.Context, prj *project.Project, logger util.Logger) ([]string, error) {
	var logs []string
	addLog := func(msg string, args ...any) {
		ret := fmt.Sprintf(msg, args...)
		logs = append(logs, ret)
	}

	addLog("building project [%s] in [%s]", prj.Key, prj.Path)

	exitCode, out, err := telemetry.RunProcessSimple(ctx, "bin/templates.sh", prj.Path, logger)
	if err != nil {
		return logs, err
	}
	addLog("templates.sh output for [" + prj.Key + "]:\n" + out)
	if exitCode != 0 {
		return logs, errors.Errorf("templates.sh failed with exit code [%d]", exitCode)
	}

	exitCode, out, err = telemetry.RunProcessSimple(ctx, "go mod tidy", prj.Path, logger)
	if err != nil {
		return logs, err
	}
	addLog("go mod output for [" + prj.Key + "]:\n" + out)
	if exitCode != 0 {
		return logs, errors.Errorf("\"go mod tidy\" failed with exit code [%d]", exitCode)
	}

	exitCode, out, err = telemetry.RunProcessSimple(ctx, "npm install", filepath.Join(prj.Path, "client"), logger)
	if err != nil {
		return logs, err
	}
	addLog("npm output for [" + prj.Key + "]:\n" + out)
	if exitCode != 0 {
		return logs, errors.Errorf("npm install failed with exit code [%d]", exitCode)
	}

	exitCode, out, err = telemetry.RunProcessSimple(ctx, "bin/build/client.sh", prj.Path, logger)
	if err != nil {
		return logs, err
	}
	addLog("client build output for [" + prj.Key + "]:\n" + out)
	if exitCode != 0 {
		return logs, errors.Errorf("client build failed with exit code [%d]", exitCode)
	}

	exitCode, out, err = telemetry.RunProcessSimple(ctx, "make build", prj.Path, logger)
	if err != nil {
		return logs, err
	}
	addLog("build output for [" + prj.Key + "]:\n" + out)
	if exitCode != 0 {
		return logs, errors.Errorf("build failed with exit code [%d]", exitCode)
	}

	return logs, nil
}
