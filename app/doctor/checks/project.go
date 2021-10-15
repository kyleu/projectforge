package checks

import (
	"fmt"

	"github.com/kyleu/projectforge/app/doctor"
	"github.com/kyleu/projectforge/app/filesystem"
	"github.com/kyleu/projectforge/app/project"
	"github.com/kyleu/projectforge/app/util"
	"go.uber.org/zap"
)

var prj = &doctor.Check{
	Key:     "project",
	Section: "app",
	Title:   "Project",
	Summary: "Verifies the Project Forge project in the working directory",
	URL:     util.AppURL,
	UsedBy:  util.AppName,
	Fn:      checkProject,
	Solve:   solveProject,
}

func checkProject(r *doctor.Result, logger *zap.SugaredLogger) *doctor.Result {
	dir := "."
	fs := filesystem.NewFileSystem(dir, logger)
	if !fs.Exists(project.ConfigFilename) {
		return r.WithError(doctor.NewError("missing", "no project found in [%s]", dir))
	}
	b, err := fs.ReadFile(project.ConfigFilename)
	if err != nil {
		return r.WithError(doctor.NewError("missing", "unable to read project from [%s]", dir))
	}
	p := &project.Project{}
	err = util.FromJSON(b, p)
	if err != nil {
		return r.WithError(doctor.NewError("invalid", "unable to parse project JSON from [%s]", dir))
	}
	if p.Port == 0 {
		r = r.WithError(doctor.NewError("config", "port must be a non-zero integer"))
	}
	if p.Info == nil {
		return r.WithError(doctor.NewError("config", "project [%s] has no project info", p.Key))
	}
	if len(p.Modules) == 0 {
		r = r.WithError(doctor.NewError("config", "No modules enabled?!"))
	}

	hasMod := func(key string) bool {
		for _, m := range p.Modules {
			if m == key {
				return true
			}
		}
		return false
	}
	if p.Build == nil {
		p.Build = &project.Build{}
	}
	if hasMod("desktop") && (!p.Build.Desktop) {
		r = r.WithError(doctor.NewError("config", "desktop module is enabled, but desktop build isn't set"))
	}
	if hasMod("ios") && (!p.Build.IOS) {
		r = r.WithError(doctor.NewError("config", "iOS module is enabled, but iOS build isn't set"))
	}
	if hasMod("android") && (!p.Build.Android) {
		r = r.WithError(doctor.NewError("config", "Android module is enabled, but Android build isn't set"))
	}
	if p.Build.Notarize && p.Info.SigningIdentity == "" {
		r = r.WithError(doctor.NewError("config", "Notarizing is enabled, but [Signing Idenitity] isn't set"))
	}
	if p.Info.Homepage == "" {
		r = r.WithError(doctor.NewError("config", "No homepage set"))
	}
	if p.Info.License == "" {
		r = r.WithError(doctor.NewError("config", "No license set"))
	}
	if p.Info.AuthorName == "" {
		r = r.WithError(doctor.NewError("config", "No author name set"))
	}
	if p.Info.AuthorEmail == "" {
		r = r.WithError(doctor.NewError("config", "No author email set"))
	}
	return r
}

func solveProject(r *doctor.Result, logger *zap.SugaredLogger) *doctor.Result {
	addSol:= func(s string) {
		if r.Solution != "" {
			r.Solution += "\n"
		}
		r.Solution += s
	}
	if r.Errors.Find("missing") != nil {
		addSol("run [projectforge create] in this directory")
	}
	if r.Errors.Find("invalid") != nil {
		addSol("the project file isn't valid JSON, not sure what you can do")
	}
	if r.Errors.Find("config") != nil {
		addSol(fmt.Sprintf("use the %s UI to configure your project", util.AppName))
	}
	return r
}
