package action

import (
	"context"
	"path/filepath"

	"github.com/pkg/errors"
	"github.com/samber/lo"

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
	return lo.FindOrElse(b, nil, func(x *Build) bool {
		return x.Key == key
	})
}

var (
	buildStart         = &Build{Key: "start", Title: "Start", Description: "Starts the prebuilt project binary", Run: onStart}
	buildDeps          = &Build{Key: "deps", Title: "Dependencies", Description: "Manages Go dependencies", Run: onDeps}
	buildImports       = &Build{Key: "imports", Title: "Imports", Description: "Reorders the imports", Run: onImports}
	buildIgnored       = &Build{Key: "ignored", Title: "Ignored", Description: "Shows files that are ignored by code generation", Run: onIgnored}
	buildPackages      = &Build{Key: "packages", Title: "Packages", Description: "Visualize your application's packages", Run: onPackages}
	buildCleanup       = &Build{Key: "cleanup", Title: "Cleanup", Description: "Cleans up file permissions", Run: onCleanup}
	buildBuild         = simpleBuild("build", "Build", "make build")
	buildClean         = simpleBuild("clean", "Clean", "make clean")
	buildTidy          = simpleBuild("tidy", "Tidy", "go mod tidy")
	buildFormat        = simpleBuild("format", "Format", "bin/format.sh")
	buildLint          = simpleBuild("lint", "Lint", "bin/check.sh")
	buildClientInstall = &Build{
		Key: "clientInstall", Title: "Client Install", Description: ciDesc, Run: func(ctx context.Context, pm *PrjAndMods, ret *Result) *Result {
			return simpleProc(ctx, "npm install", filepath.Join(pm.Prj.Path, "client"), ret, pm.Logger)
		},
	}
	buildClientBuild = simpleBuild("clientBuild", "Client Build", "bin/build/client.sh")
	buildDeployments = &Build{Key: "deployments", Title: "Deployments", Description: "Manages deployments", Run: onDeployments}
	buildTest        = &Build{Key: "test", Title: "Test", Description: "Does a test", Run: func(ctx context.Context, pm *PrjAndMods, ret *Result) *Result {
		return simpleProc(ctx, "bin/test.sh", pm.Prj.Path, ret, pm.Logger)
	}}
)

var AllBuilds = Builds{
	buildStart, buildDeps, buildImports, buildIgnored, buildPackages, buildCleanup, buildBuild, buildClean,
	buildTidy, buildFormat, buildLint, buildClientInstall, buildClientBuild, buildDeployments, buildTest,
}

func fullBuild(ctx context.Context, prj *project.Project, r *Result, logger util.Logger) *Result {
	logs, err := build.Full(ctx, prj, logger)
	lo.ForEach(logs, func(l string, _ int) {
		r.AddLog(l)
	})
	if err != nil {
		return r.WithError(err)
	}
	return r
}
