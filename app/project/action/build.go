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

func (b Builds) ForAllProjects() Builds {
	return lo.Filter(b, func(x *Build, _ int) bool {
		return x.Key != "lint"
	})
}

func (b Builds) Get(key string) *Build {
	return lo.FindOrElse(b, nil, func(x *Build) bool {
		return x.Key == key
	})
}

var (
	buildFull          = &Build{Key: "full", Title: "Full Build", Description: "Builds the TypeScript and Go code", Run: onFull}
	buildBuild         = simpleBuild("build", "Build", "make build")
	buildStart         = &Build{Key: "start", Title: "Start", Description: "Starts the prebuilt project binary", Run: onStart}
	buildClean         = simpleBuild("clean", "Clean", "make clean")
	buildDeps          = &Build{Key: "deps", Title: "Dependencies", Description: "Manages Go dependencies", Run: onDeps}
	buildImports       = &Build{Key: "imports", Title: "Imports", Description: "Reorders the imports", Run: onImports}
	buildIgnored       = &Build{Key: "ignored", Title: "Ignored", Description: "Shows files that are ignored by code generation", Run: onIgnored}
	buildPackages      = &Build{Key: "packages", Title: "Packages", Description: "Visualize your application's packages", Run: onPackages}
	buildCleanup       = &Build{Key: "cleanup", Title: "Cleanup", Description: "Cleans up file permissions", Run: onCleanup}
	buildTidy          = simpleBuild("tidy", "Tidy", "go mod tidy")
	buildFormat        = simpleBuild("format", "Format", filepath.Join("bin", "format."+build.ScriptExtension))
	buildLint          = simpleBuild("lint", "Lint", filepath.Join("bin", "check."+build.ScriptExtension))
	buildLintClient    = simpleBuild("lint-client", "Lint Client", filepath.Join("bin", "check-client."+build.ScriptExtension))
	buildTemplates     = simpleBuild("templates", "Templates", filepath.Join("bin", "templates."+build.ScriptExtension))
	buildClientInstall = &Build{
		Key: "clientInstall", Title: "Client Install", Description: ciDesc, Run: func(ctx context.Context, pm *PrjAndMods, ret *Result) *Result {
			return simpleProc(ctx, "npm install", filepath.Join(pm.Prj.Path, "client"), ret, pm.Logger)
		},
	}
	buildClientBuild = simpleBuild("clientBuild", "Client Build", filepath.Join("bin", "build", "client."+build.ScriptExtension))
	buildDeployments = &Build{Key: "deployments", Title: "Deployments", Description: "Manages deployments", Run: onDeployments}
	buildTest        = &Build{Key: "test", Title: "Test", Description: "Runs unit tests", Run: func(ctx context.Context, pm *PrjAndMods, ret *Result) *Result {
		return simpleProc(ctx, filepath.Join("bin", "test."+build.ScriptExtension), pm.Prj.Path, ret, pm.Logger)
	}}
)

var AllBuilds = Builds{
	buildFull, buildBuild, buildStart, buildClean, buildDeps, buildImports, buildIgnored, buildPackages, buildCleanup,
	buildTidy, buildFormat, buildLint, buildLintClient, buildTemplates, buildClientInstall, buildClientBuild, buildDeployments, buildTest,
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
