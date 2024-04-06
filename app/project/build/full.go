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

func Full(ctx context.Context, prj *project.Project, logger util.Logger) ([]string, error) {
	var logs []string
	addLog := func(msg string, args ...any) {
		ret := fmt.Sprintf(msg, args...)
		logs = append(logs, ret)
	}
	addLogOutput := func(key string, out string) {
		addLog("%s output for [%s]:\n%s", key, prj.Key, out)
	}
	cmd := func(key string, cmd string, pth string) error {
		if key == "" {
			key = cmd
		}
		if pth == "" {
			pth = prj.Path
		}
		exitCode, out, err := telemetry.RunProcessSimple(ctx, cmd, pth, logger)
		if err != nil {
			return err
		}
		addLogOutput(key, out)
		if exitCode != 0 {
			return errors.Errorf(key+" failed with exit code [%d]", exitCode)
		}
		return nil
	}

	addLog("building project [%s] in [%s]", prj.Key, prj.Path)
	err := cmd(templatesNS+ScriptExtension, filepath.Join("bin", templatesNS+ScriptExtension), "")
	if err != nil {
		return logs, err
	}
	err = cmd("", "go mod tidy", "")
	if err != nil {
		return logs, err
	}
	err = cmd("", "npm install", filepath.Join(prj.Path, "client"))
	if err != nil {
		return logs, err
	}
	err = cmd("client build", filepath.Join("bin", "build", "client."+ScriptExtension), "")
	if err != nil {
		return logs, err
	}
	makeCmd := "make build"
	if runtime.GOOS == OSWindows {
		makeCmd = fmt.Sprintf(`go build -ldflags "-s -w" -trimpath -o build/release/%s.exe`, prj.Executable())
	}
	err = cmd("project build", makeCmd, "")
	if err != nil {
		return logs, err
	}
	return logs, nil
}
