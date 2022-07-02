package action

import (
	"context"
	"path/filepath"

	"github.com/pkg/errors"

	"projectforge.dev/projectforge/app/lib/telemetry"
	"projectforge.dev/projectforge/app/project"
	"projectforge.dev/projectforge/app/project/build"
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
	{Key: "packages", Title: "Packages", Description: "Visualize your application's packages", Run: onPackages},
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
