package project

import (
	"fmt"
	"slices"
	"strings"

	"github.com/samber/lo"

	"projectforge.dev/projectforge/app/lib/filesystem"
)

const DebugOutputDir = "build/debug/"

type validationAddErrFn func(code string, msg string, args ...any)

type ValidationError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

func Validate(p *Project, fs filesystem.FileLoader, moduleDeps map[string][]string, dangerous []string) []*ValidationError {
	var ret []*ValidationError

	e := func(code string, msg string, args ...any) {
		if len(args) > 0 {
			msg = fmt.Sprintf(msg, args...)
		}
		ret = append(ret, &ValidationError{Code: code, Message: msg})
	}

	validateBasic(p, e)
	validateFilesystem(p, e, fs)
	validateModuleDeps(p.Modules, moduleDeps, e)
	validateModuleConfig(p, e, dangerous)
	validateModuleServices(p.Modules, fs, e)
	validateInfo(p, e)
	validateBuild(p, e)
	validateExport(p, e)

	return ret
}

func validateBasic(p *Project, e validationAddErrFn) {
	if p.Port == 0 {
		e("port", "port must be a non-zero integer")
	}
	if p.Info == nil {
		e("nil-info", "no project info")
		p.Info = &Info{}
	}
	if len(p.Modules) == 0 {
		e("no-modules", "no modules enabled")
	}
	if !lo.Contains(p.Modules, "core") {
		e("no-modules", "core module not included")
	}
}

func validateModuleDeps(modules []string, deps map[string][]string, e validationAddErrFn) {
	if deps == nil {
		return
	}
	lo.ForEach(modules, func(m string, _ int) {
		if currDeps, ok := deps[m]; ok && len(currDeps) > 0 {
			lo.ForEach(currDeps, func(curr string, _ int) {
				if !lo.Contains(modules, curr) {
					e("missing-dependency", "module [%s] requires [%s], which is not included in the project", m, curr)
				}
			})
		}
	})
}

func validateModuleServices(modules []string, fs filesystem.FileLoader, e validationAddErrFn) {
	b, err := fs.ReadFile("app/services.go")
	if err != nil {
		e("missing-file", "missing [app/services.go]")
	}

	servicesContent := string(b)
	lo.ForEach(modules, func(m string, _ int) {
		if key := templateServicesKeys[m]; key != "" && key != "user" {
			svc := fmt.Sprintf("*%s.Service", key)
			if !strings.Contains(servicesContent, svc) {
				e("missing-services-reference", "module [%s] requires [%s], which is not included in [app/services.go]", m, svc)
			}
		}
	})
}

func validateModuleConfig(p *Project, e validationAddErrFn, dangerous []string) {
	if p.HasModule("desktop") && (!p.Build.Desktop) {
		e("desktop-disabled", "desktop module is enabled, but desktop build isn't set")
	}
	if p.HasModule("ios") && (!p.Build.IOS) {
		e("ios-disabled", "iOS module is enabled, but iOS build isn't set")
	}
	if p.HasModule("android") && (!p.Build.Android) {
		e("android-disabled", "Android module is enabled, but Android build isn't set")
	}
	if p.Build.SafeMode {
		lo.ForEach(dangerous, func(m string, _ int) {
			if p.HasModule(m) {
				e("dangerous-module", "Safe mode is enabled for this project, but dangerous module [%s] is enabled", m)
			}
		})
	}
}

func validateBuild(p *Project, e validationAddErrFn) {
	if p.Build == nil {
		p.Build = &Build{}
	}

	if p.Build.Desktop && !lo.Contains(p.Modules, "desktop") {
		e("config", "Desktop is enabled, but module [desktop] isn't included")
	}
	if p.Build.Desktop && lo.Contains(p.Modules, "desktop") && p.Info.Bundle == "" {
		e("config", "Desktop build is enabled, but [Bundle] isn't set")
	}

	if p.Build.IOS && !lo.Contains(p.Modules, "ios") {
		e("config", "iOS build is enabled, but module [ios] isn't included")
	}
	if p.Build.IOS && lo.Contains(p.Modules, "ios") && p.Info.Bundle == "" {
		e("config", "iOS build is enabled, but [Bundle] isn't set")
	}

	if p.Build.Android && !lo.Contains(p.Modules, "android") {
		e("config", "Android build is enabled, but module [android] isn't included")
	}
	if p.Build.Android && lo.Contains(p.Modules, "android") && p.Info.JavaPackage == "" {
		e("config", "Android build is enabled, but [Java Package] isn't set")
	}

	if p.Build.Notarize && !lo.Contains(p.Modules, "notarize") {
		e("config", "Notarize build is enabled, but module [notarize] isn't included")
	}
	if p.Build.Notarize && p.Info.SigningIdentity == "" {
		e("config", "Notarizing is enabled, but [Signing Identity] isn't set")
	}
	if p.Build.Notarize && p.Info.Bundle == "" {
		e("config", "Notarizing is enabled, but [Bundle] isn't set")
	}

	if p.Build.WASM && !lo.Contains(p.Modules, "wasmserver") {
		e("config", "WASM build is enabled, but module [wasmserver] isn't included")
	}
}

func validateInfo(p *Project, e validationAddErrFn) {
	if p.Info.Homepage == "" {
		e("config", "No homepage set")
	}
	if p.Info.License == "" {
		e("config", "No license set")
	}
	if p.Info.AuthorID == "" {
		e("config", "No author ID set")
	}
	if p.Info.AuthorName == "" {
		e("config", "No author name set")
	}
	if p.Info.AuthorEmail == "" {
		e("config", "No author email set")
	}
}

func validateExport(p *Project, e validationAddErrFn) {
	if p.ExportArgs == nil {
		return
	}
	if err := p.ExportArgs.Validate(); err != nil {
		e("export", err.Error())
	}
	if err := p.ExportArgs.Models.Validate(p.Modules, p.ExportArgs.Groups); err != nil {
		e("export", err.Error())
	}
}

func validateFilesystem(p *Project, e validationAddErrFn, fs filesystem.FileLoader) {
	if fs == nil {
		e("missing-filesystem", "The project filesystem does not exist")
		return
	}
	if !fs.Exists(".projectforge/project.json") {
		e("project-file", "the project definition file (.projectforge/project.json) is missing")
	}
	if !fs.Exists("app/services.go") {
		e("needs-generate", "some generated files are missing, run the \"Generate\" action")
		return
	}

	if !fs.Exists("go.mod") {
		e("missing-files", "this project needs to be generated using the button above, or using the CLI")
		return
	}
	if !fs.Exists("views/Home.html.go") {
		e("no-template-build", "it looks like your templates haven't been generated, perform a build or run \"bin/templates.sh\"")
		return
	}
	if !fs.Exists("go.sum") {
		e("go-dependencies", "the Go dependencies file is missing, run \"go mod tidy\" or perform a full build")
		return
	}
	if !fs.Exists("client/package-lock.json") {
		e("ts-dependencies", "the TypeScript dependencies file is missing, run \"npm i\" or perform a full build")
		return
	}
	if !fs.Exists("assets/client.css") {
		e("front-end-build", "the TypeScript build output is missing, run \"bin/build/client.sh\" or perform a full build")
		return
	}
	if !fs.Exists(DebugOutputDir+p.Executable()) && !fs.Exists(DebugOutputDir+p.Executable()+".exe") {
		e("needs-build", "your project hasn't been built recently, run a build using the buttons above")
		return
	}

	if slices.Contains(p.Modules, "export") {
		if !fs.Exists(".projectforge/export") {
			e("missing-export-directory", "the project uses the export module, but doesn't have directory [./projectforge/export]")
		}
	}

	if b, err := fs.ReadFile("client/src/svg/app.svg"); err == nil {
		if strings.Contains(string(b), "default_icon") {
			e("default-icon", "this project uses the default application icon, choose a new one using the SVG manager")
		}
	} else {
		e("no-icon", "the file that defines the main icon (client/src/svg/app.svg) is missing")
	}
}
