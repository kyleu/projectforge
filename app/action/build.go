package action

import (
	"context"
	"path/filepath"

	"github.com/pkg/errors"

	"projectforge.dev/projectforge/app/build"
	"projectforge.dev/projectforge/app/lib/telemetry"
	"projectforge.dev/projectforge/app/module"
	"projectforge.dev/projectforge/app/project"
	"projectforge.dev/projectforge/app/util"
)

type Build struct {
	Key         string `json:"key"`
	Title       string `json:"title,omitempty"`
	Description string `json:"description,omitempty"`

	Run func(ctx context.Context, pm *PrjAndMods, ret *Result) *Result `json:"-"`
}

func simpleProc(ctx context.Context, cmd string, path string, ret *Result, logger util.Logger) *Result {
	exitCode, out, err := telemetry.RunProcessSimple(ctx, cmd, path, logger)
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
	return &Build{Key: key, Title: title, Description: "Runs [" + cmd + "]", Run: func(ctx context.Context, pm *PrjAndMods, ret *Result) *Result {
		return simpleProc(ctx, cmd, pm.Prj.Path, ret, pm.Logger)
	}}
}

const ciDesc = "Installs dependencies for the TypeScript client"

type Builds []*Build

func (b Builds) Get(key string) *Build {
	for _, x := range b {
		if x.Key == key {
			return x
		}
	}
	return nil
}

var AllBuilds = Builds{
	{Key: "deps", Title: "Dependencies", Description: "Manages Go dependencies", Run: onDeps},
	{Key: "imports", Title: "Imports", Description: "Reorders the imports", Run: onImports},
	{Key: "cleanup", Title: "Cleanup", Description: "Cleans up file permissions", Run: onCleanup},
	simpleBuild("build", "Build", "make build"),
	simpleBuild("clean", "Clean", "make clean"),
	simpleBuild("tidy", "Tidy", "go mod tidy"),
	simpleBuild("format", "Format", "bin/format.sh"),
	simpleBuild("lint", "Lint", "bin/check.sh"),
	{Key: "clientInstall", Title: "Client Install", Description: ciDesc, Run: func(ctx context.Context, pm *PrjAndMods, ret *Result) *Result {
		return simpleProc(ctx, "npm install", filepath.Join(pm.Prj.Path, "client"), ret, pm.Logger)
	}},
	simpleBuild("clientBuild", "Client Build", "bin/build/client.sh"),
	{Key: "deployments", Title: "Deployments", Description: "Manages deployments", Run: onDeployments},
	{Key: "test", Title: "Test", Description: "Does a test", Run: func(ctx context.Context, pm *PrjAndMods, ret *Result) *Result {
		return simpleProc(ctx, "ls", pm.Prj.Path, ret, pm.Logger)
	}},
}

func onBuild(ctx context.Context, pm *PrjAndMods) *Result {
	phaseStr := pm.Cfg.GetStringOpt("phase")
	if phaseStr == "" {
		phaseStr = "build"
	}

	_, span, logger := telemetry.StartSpan(ctx, "build:"+phaseStr, pm.Logger)
	defer span.Complete()
	pm.Logger = logger

	ret := newResult(TypeBuild, pm.Prj, pm.Cfg, logger)
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
	timer := util.TimerStart()
	ret = phase.Run(ctx, pm, ret)
	ret.Duration = timer.End()
	return ret
}

func onDeps(ctx context.Context, pm *PrjAndMods, ret *Result) *Result {
	if up := ret.Args.GetStringOpt("upgrade"); up != "" {
		o := ret.Args.GetStringOpt("old")
		n := ret.Args.GetStringOpt("new")
		err := build.OnDepsUpgrade(ctx, pm.Prj, up, o, n, pm.PSvc, pm.Logger)
		if err != nil {
			return ret.WithError(err)
		}
	}
	upd := ret.Args.GetStringOpt("upd")
	showAll := pm.Cfg.GetBoolOpt("showAll")
	deps, err := build.LoadDeps(ctx, pm.Prj.Key, pm.Prj.Path, upd != "false", pm.PSvc.GetFilesystem(pm.Prj), showAll, pm.Logger)
	ret.Data = deps
	if err != nil {
		return ret.WithError(err)
	}
	return ret
}

func onImports(ctx context.Context, pm *PrjAndMods, r *Result) *Result {
	fixStr := r.Args.GetStringOpt("fix")
	fix := fixStr == "true"
	fileStr := r.Args.GetStringOpt("file")
	t := util.TimerStart()
	logs, diffs, err := build.Imports(pm.Prj.Package, fix, fileStr, pm.PSvc.GetFilesystem(pm.Prj), pm.Logger)
	r.Modules = append(r.Modules, &module.Result{Keys: []string{"imports"}, Status: "OK", Diffs: diffs, Duration: t.End()})
	r.Logs = append(r.Logs, logs...)
	if err != nil {
		return r.WithError(err)
	}
	return r
}

func onCleanup(ctx context.Context, pm *PrjAndMods, r *Result) *Result {
	t := util.TimerStart()
	logs, diffs, err := build.Cleanup(pm.PSvc.GetFilesystem(pm.Prj), pm.Logger)
	r.Modules = append(r.Modules, &module.Result{Keys: []string{"cleanup"}, Status: "OK", Diffs: diffs, Duration: t.End()})
	r.Logs = append(r.Logs, logs...)
	if err != nil {
		return r.WithError(err)
	}
	return r
}

func onDeployments(ctx context.Context, pm *PrjAndMods, r *Result) *Result {
	if pm.Prj.Info == nil || len(pm.Prj.Info.Deployments) == 0 {
		return r
	}
	t := util.TimerStart()
	fix := pm.Cfg.GetBoolOpt("fix")
	logs, diffs, err := build.Deployments(ctx, pm.Prj.Version, pm.PSvc.GetFilesystem(pm.Prj), fix, pm.Cfg.GetStringOpt("path"), pm.Prj.Info.Deployments)
	r.Modules = append(r.Modules, &module.Result{Keys: []string{"deployments"}, Status: "OK", Diffs: diffs, Duration: t.End()})
	r.Logs = append(r.Logs, logs...)
	if err != nil {
		return r.WithError(err)
	}
	return r
}

func fullBuild(ctx context.Context, prj *project.Project, r *Result, logger util.Logger) *Result {
	logs, err := build.Full(ctx, prj, logger)
	for _, l := range logs {
		r.AddLog(l)
	}
	if err != nil {
		return r.WithError(err)
	}
	return r
}
