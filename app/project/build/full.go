package build

import (
	"context"
	"fmt"
	"path/filepath"
	"runtime"

	"github.com/pkg/errors"

	"projectforge.dev/projectforge/app/lib/telemetry"
	"projectforge.dev/projectforge/app/project"
	"projectforge.dev/projectforge/app/util"
)

const templatesNS = "templates."

type ExecHelper struct {
	Logs []string `json:"logs,omitempty"`
}

func (e *ExecHelper) AddLog(msg string, args ...any) {
	e.Logs = append(e.Logs, fmt.Sprintf(msg, args...))
}

func (e *ExecHelper) AddLogOutput(key string, out string) {
	e.AddLog("%s output: %s", key, out)
}

func (e *ExecHelper) Cmd(ctx context.Context, key string, cmd string, pth string, logger util.Logger) (string, error) {
	if key == "" {
		key = cmd
	}
	exitCode, out, err := telemetry.RunProcessSimple(ctx, cmd, pth, logger)
	if err != nil {
		return "", err
	}
	e.AddLogOutput(key, out)
	if exitCode != 0 {
		return out, errors.Errorf(key+" failed with exit code [%d]", exitCode)
	}
	return out, nil
}

func Full(ctx context.Context, prj *project.Project, logger util.Logger) ([]string, error) {
	ex := &ExecHelper{}
	ex.AddLog("building project [%s] in [%s]", prj.Key, prj.Path)
	_, err := ex.Cmd(ctx, templatesNS+ScriptExtension, filepath.Join("bin", templatesNS+ScriptExtension), "", logger)
	if err != nil {
		return ex.Logs, err
	}
	_, err = ex.Cmd(ctx, "", "go mod tidy", "", logger)
	if err != nil {
		return ex.Logs, err
	}
	_, err = ex.Cmd(ctx, "", "npm install", filepath.Join(prj.Path, "client"), logger)
	if err != nil {
		return ex.Logs, err
	}
	_, err = ex.Cmd(ctx, "client build", filepath.Join("bin", "build", "client."+ScriptExtension), "", logger)
	if err != nil {
		return ex.Logs, err
	}
	makeCmd := "make build"
	if runtime.GOOS == OSWindows {
		makeCmd = fmt.Sprintf(`go build -ldflags "-s -w" -trimpath -o build/release/%s.exe`, prj.Executable())
	}
	_, err = ex.Cmd(ctx, "project build", makeCmd, "", logger)
	if err != nil {
		return ex.Logs, err
	}
	return ex.Logs, nil
}
