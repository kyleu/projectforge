package build

import (
	"fmt"
	"path/filepath"

	"github.com/pkg/errors"
	"go.uber.org/zap"
	"projectforge.dev/projectforge/app/project"
	"projectforge.dev/projectforge/app/util"
)

func Full(prj *project.Project, logger *zap.SugaredLogger) ([]string, error) {
	var logs []string
	addLog := func(msg string, args ...any) {
		ret := fmt.Sprintf(msg, args...)
		logs = append(logs, ret)
	}

	addLog("building project [%s] in [%s]", prj.Key, prj.Path)

	exitCode, out, err := util.RunProcessSimple("bin/templates.sh", prj.Path)
	if err != nil {
		return logs, err
	}
	addLog("templates.sh output for [" + prj.Key + "]:\n" + out)
	if exitCode != 0 {
		return logs, errors.Errorf("templates.sh failed with exit code [%d]", exitCode)
	}

	exitCode, out, err = util.RunProcessSimple("go mod tidy", prj.Path)
	if err != nil {
		return logs, err
	}
	addLog("go mod output for [" + prj.Key + "]:\n" + out)
	if exitCode != 0 {
		return logs, errors.Errorf("\"go mod tidy\" failed with exit code [%d]", exitCode)
	}

	exitCode, out, err = util.RunProcessSimple("npm install", filepath.Join(prj.Path, "client"))
	if err != nil {
		return logs, err
	}
	addLog("npm output for [" + prj.Key + "]:\n" + out)
	if exitCode != 0 {
		return logs, errors.Errorf("npm install failed with exit code [%d]", exitCode)
	}

	exitCode, out, err = util.RunProcessSimple("bin/build/client.sh", prj.Path)
	if err != nil {
		return logs, err
	}
	addLog("client build output for [" + prj.Key + "]:\n" + out)
	if exitCode != 0 {
		return logs, errors.Errorf("client build failed with exit code [%d]", exitCode)
	}

	exitCode, out, err = util.RunProcessSimple("make build", prj.Path)
	if err != nil {
		return logs, err
	}
	addLog("build output for [" + prj.Key + "]:\n" + out)
	if exitCode != 0 {
		return logs, errors.Errorf("build failed with exit code [%d]", exitCode)
	}

	return logs, nil
}
