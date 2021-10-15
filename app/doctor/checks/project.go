package checks

import (
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
	return r
}
