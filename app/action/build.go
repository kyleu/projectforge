package action

import (
	"path/filepath"

	"github.com/kyleu/projectforge/app/util"
	"github.com/pkg/errors"
)

type Build struct {
	Key         string `json:"key"`
	Title       string `json:"title"`
	Description string `json:"description"`

	Run func(pm *PrjAndMods, ret *Result) *Result `json:"-"`
}

func simpleProc(cmd string, path string, ret *Result) *Result {
	exitCode, out, err := util.RunProcessSimple(cmd, path)
	if err != nil {
		return ret.WithError(err)
	}
	ret.AddLog("build output for [" + cmd + "]:\n" + out)
	if exitCode != 0 {
		ret.WithError(errors.Errorf("build failed with exit code [%d]", exitCode))
	}
	return ret
}

func simpleBuild(key string, title string, cmd string) *Build {
	return &Build{Key: key, Title: title, Description: "Runs [" + cmd + "]", Run: func(pm *PrjAndMods, ret *Result) *Result {
		return simpleProc(cmd, pm.Prj.Path, ret)
	}}
}

const ciDesc = "Installs dependencies for the TypeScript client"

var AllBuilds = []*Build{
	simpleBuild("build", "Build", "make build"),
	simpleBuild("clean", "Clean", "make clean"),
	simpleBuild("tidy", "Tidy", "go mod tidy"),
	simpleBuild("format", "Format", "bin/format.sh"),
	simpleBuild("lint", "Lint", "bin/check.sh"),
	{Key: "clientInstall", Title: "Client Install", Description: ciDesc, Run: func(pm *PrjAndMods, ret *Result) *Result {
		return simpleProc("npm install", filepath.Join(pm.Prj.Path, "client"), ret)
	}},
	simpleBuild("clientBuild", "Client Build", "bin/build/client.sh"),
	{Key: "test", Title: "Test", Description: "Does a test", Run: func(pm *PrjAndMods, ret *Result) *Result {
		return simpleProc("ls", pm.Prj.Path, ret)
	}},
}

func onBuild(pm *PrjAndMods, cfg util.ValueMap) *Result {
	phaseStr, _ := cfg.GetString("phase", true)
	if phaseStr == "" {
		phaseStr = "make"
	}

	ret := newResult(pm.Cfg, pm.Logger)
	ret.AddLog("building project [%s] in [%s] with phase [%s]", pm.Prj.Key, pm.Prj.Path, phaseStr)
	var phase *Build
	for _, x := range AllBuilds {
		if x.Key == phaseStr {
			phase = x
			break
		}
	}
	if phase == nil {
		return ret.WithError(errors.Errorf("invalid phase [%s]", phaseStr))
	}
	start := util.TimerStart()
	ret = phase.Run(pm, ret)
	ret.Duration = util.TimerEnd(start)
	return ret
}
