package action

import (
	"context"

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
	Expensive   bool   `json:"expensive,omitempty"`

	Run func(ctx context.Context, pm *PrjAndMods, ret *Result) *Result `json:"-"`
}

func simpleProc(ctx context.Context, cmd string, path string, ret *Result, logger util.Logger) *Result {
	exitCode, out, err := telemetry.RunProcessSimple(ctx, cmd, path, logger)
	if err != nil {
		return ret.WithError(err)
	}
	ret.AddLog("build output for [%s]:\n%s", cmd, out)
	if exitCode != 0 {
		ret.WithError(errors.Errorf("build failed with exit code [%d]", exitCode))
	}
	return ret
}

func simpleBuild(key string, title string, cmd string, expensive bool) *Build {
	return &Build{Key: key, Title: title, Description: "Runs [" + cmd + "]", Run: func(ctx context.Context, pm *PrjAndMods, ret *Result) *Result {
		return simpleProc(ctx, cmd, pm.Prj.Path, ret, pm.Logger)
	}, Expensive: expensive}
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
	buildBuild         = simpleBuild("build", "Build", "make build", false)
	buildStart         = &Build{Key: "start", Title: "Start", Description: "Starts the prebuilt project binary", Run: onStart}
	buildClean         = simpleBuild("clean", "Clean", "make clean", false)
	buildDeps          = &Build{Key: "deps", Title: "Dependencies", Description: "Manages Go dependencies", Run: onDeps}
	buildImports       = &Build{Key: "imports", Title: "Imports", Description: "Organizes the imports in source files and templates", Run: onImports}
	buildIgnored       = &Build{Key: "ignored", Title: "Ignored", Description: "Shows files that are ignored by code generation", Run: onIgnored}
	buildPackages      = &Build{Key: "packages", Title: "Packages", Description: "Visualize your application's packages", Run: onPackages}
	buildCleanup       = &Build{Key: "cleanup", Title: "Cleanup", Description: "Cleans up file permissions", Run: onCleanup}
	buildSize          = &Build{Key: "size", Title: "Binary Size", Description: "Visualizes the file size of the binary", Run: onSize}
	buildTidy          = simpleBuild("tidy", "Tidy", "go mod tidy", false)
	buildFormat        = simpleBuild("format", "Format", util.StringFilePath("bin", "format."+build.ScriptExtension), false)
	buildFormatClient  = simpleBuild("format-client", "Format Client", util.StringFilePath("bin", "format-client."+build.ScriptExtension), false)
	buildLint          = simpleBuild("lint", "Lint", util.StringFilePath("bin", "check."+build.ScriptExtension), true)
	buildLintClient    = simpleBuild("lint-client", "Lint Client", util.StringFilePath("bin", "check-client."+build.ScriptExtension), false)
	buildTemplates     = simpleBuild("templates", "Templates", util.StringFilePath("bin", "templates."+build.ScriptExtension), false)
	buildClientInstall = &Build{Key: "clientInstall", Title: "Client Install", Description: ciDesc, Run: onClientInstall}
	buildClientBuild   = simpleBuild("clientBuild", "Client Build", util.StringFilePath("bin", "build", "client."+build.ScriptExtension), false)
	buildThemeRebuild  = &Build{Key: "themeRebuild", Title: "Theme Rebuild", Description: "Rebuilds the theme", Run: onThemeRebuild}
	buildDeployments   = &Build{Key: "deployments", Title: "Deployments", Description: "Manages deployments", Run: onDeployments}
	buildCoverage      = &Build{Key: "coverage", Title: "Code Coverage", Description: "Runs unit tests, displaying a coverage report", Run: onCoverage, Expensive: true}
	buildTest          = &Build{Key: "test", Title: "Test", Description: "Runs unit tests", Run: func(ctx context.Context, pm *PrjAndMods, ret *Result) *Result {
		return simpleProc(ctx, util.StringFilePath("bin", "test."+build.ScriptExtension), pm.Prj.Path, ret, pm.Logger)
	}, Expensive: true}
)

var AllBuilds = Builds{
	buildFull, buildBuild, buildStart, buildClean, buildDeps, buildImports, buildIgnored, buildPackages, buildCleanup, buildSize,
	buildTidy, buildFormat, buildFormatClient, buildLint, buildLintClient, buildTemplates, buildClientInstall, buildClientBuild,
	buildThemeRebuild, buildDeployments, buildTest, buildCoverage,
}

func fullBuild(ctx context.Context, prj *project.Project, r *Result, logger util.Logger) *Result {
	logs, err := build.Full(ctx, prj, logger)
	lo.ForEach(logs, func(l string, _ int) {
		r.Log(l)
	})
	if err != nil {
		return r.WithError(err)
	}
	return r
}
