package action

import (
	"projectforge.dev/projectforge/app/project/build"
	"projectforge.dev/projectforge/app/util"
)

const (
	ciDesc = "Installs dependencies for the TypeScript client"
	cuDesc = "Updates dependency versions for the TypeScript client"
)

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
	buildFix           = simpleBuild("fix", "Go Fix", "go fix ./...", true)
	buildFormat        = simpleBuild("format", "Format", util.StringFilePath("bin", "format."+build.ScriptExtension), false)
	buildFormatClient  = simpleBuild("format-client", "Format Client", util.StringFilePath("bin", "format-client."+build.ScriptExtension), false)
	buildLint          = simpleBuild("lint", "Lint", util.StringFilePath("bin", "check."+build.ScriptExtension), true)
	buildLintClient    = simpleBuild("lint-client", "Lint Client", util.StringFilePath("bin", "check-client."+build.ScriptExtension), false)
	buildTemplates     = simpleBuild("templates", "Templates", util.StringFilePath("bin", "templates."+build.ScriptExtension), false)
	buildClientInstall = &Build{Key: "clientInstall", Title: "Client Install", Description: ciDesc, Run: onClientInstall}
	buildClientUpdate  = &Build{Key: "clientUpdate", Title: "Client Update", Description: cuDesc, Run: onClientUpdate}
	buildClientBuild   = simpleBuild("clientBuild", "Client Build", util.StringFilePath("bin", "build", "client."+build.ScriptExtension), false)
	buildThemeRebuild  = &Build{Key: "themeRebuild", Title: "Theme Rebuild", Description: "Rebuilds the theme", Run: onThemeRebuild}
	buildDeployments   = &Build{Key: "deployments", Title: "Deployments", Description: "Manages deployments", Run: onDeployments}
	buildTest          = &Build{Key: "test", Title: "Test", Description: "Runs unit tests", Run: onBuildTest, Expensive: true}
	buildTestClient    = &Build{Key: "test-client", Title: "Test Client", Description: "Runs client unit tests", Run: onBuildClientTest, Expensive: true}
	buildCoverage      = &Build{Key: "coverage", Title: "Code Coverage", Description: "Runs a coverage report", Run: onCoverage, Expensive: true}
	buildCustom        = &Build{Key: "custom", Title: "Custom Command", Description: "Runs a custom command", Run: onCustom, Expensive: true}
)

var AllBuilds = Builds{
	buildFull, buildBuild, buildStart, buildClean, buildDeps, buildImports, buildIgnored, buildPackages, buildCleanup, buildSize,
	buildTidy, buildFix, buildFormat, buildFormatClient, buildLint, buildLintClient, buildTemplates, buildClientInstall, buildClientUpdate,
	buildClientBuild, buildThemeRebuild, buildDeployments, buildTest, buildTestClient, buildCoverage, buildCustom,
}
