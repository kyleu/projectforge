package project

import (
	"fmt"

	"golang.org/x/exp/slices"

	"projectforge.dev/projectforge/app/lib/filesystem"
)

type ValidationError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

func Validate(p *Project, moduleDeps map[string][]string, fs filesystem.FileLoader) []*ValidationError {
	var ret []*ValidationError

	e := func(code string, msg string, args ...any) {
		if len(args) > 0 {
			msg = fmt.Sprintf(msg, args...)
		}
		ret = append(ret, &ValidationError{Code: code, Message: msg})
	}

	validateBasic(p, e)
	validateModuleDeps(p.Modules, moduleDeps, e)
	validateModuleConfig(p, e)
	validateBuild(p, e)
	validateExport(p, e)

	return ret
}

func validateBasic(p *Project, e func(code string, msg string, args ...any)) {
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
	if !slices.Contains(p.Modules, "core") {
		e("no-modules", "core module not included")
	}
}

func validateModuleDeps(modules []string, deps map[string][]string, e func(code string, msg string, args ...any)) {
	if deps == nil {
		return
	}
	for _, m := range modules {
		if currDeps, ok := deps[m]; ok && len(currDeps) > 0 {
			for _, curr := range currDeps {
				if !slices.Contains(modules, curr) {
					e("missing-dependency", "module [%s] requires [%s], which is not included in the project", m, curr)
				}
			}
		}
	}
}

func validateModuleConfig(p *Project, e func(code string, msg string, args ...any)) {
	if p.HasModule("desktop") && (!p.Build.Desktop) {
		e("desktop-disabled", "desktop module is enabled, but desktop build isn't set")
	}
	if p.HasModule("ios") && (!p.Build.IOS) {
		e("ios-disabled", "iOS module is enabled, but iOS build isn't set")
	}
	if p.HasModule("android") && (!p.Build.Android) {
		e("android-disabled", "Android module is enabled, but Android build isn't set")
	}
}

func validateBuild(p *Project, e func(code string, msg string, args ...any)) {
	if p.Build == nil {
		p.Build = &Build{}
	}

	if p.Build.Desktop && !slices.Contains(p.Modules, "desktop") {
		e("config", "Desktop is enabled, but module [desktop] isn't included")
	}
	if p.Build.IOS && !slices.Contains(p.Modules, "ios") {
		e("config", "IOS is enabled, but module [ios] isn't included")
	}
	if p.Build.Android && !slices.Contains(p.Modules, "android") {
		e("config", "Android is enabled, but module [android] isn't included")
	}

	if p.Build.Notarize && p.Info.SigningIdentity == "" {
		e("config", "Notarizing is enabled, but [Signing Identity] isn't set")
	}
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

func validateExport(p *Project, e func(code string, msg string, args ...any)) {
	if p.ExportArgs == nil {
		return
	}
	for _, x := range p.ExportArgs.Models {
		if err := x.Validate(p.Modules, p.ExportArgs.Groups); err != nil {
			e("export", err.Error())
		}
	}
}
