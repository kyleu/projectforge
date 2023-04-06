package checks

import (
	"context"
	"fmt"

	"projectforge.dev/projectforge/app/doctor"
	"projectforge.dev/projectforge/app/project"
	"projectforge.dev/projectforge/app/util"
)

var CurrentModuleDeps map[string][]string

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

func checkProject(_ context.Context, r *doctor.Result, _ util.Logger) *doctor.Result {
	p, _, r := loadRootProject(r)
	if len(r.Errors) > 0 {
		return r
	}
	errs := project.Validate(p, CurrentModuleDeps)
	for _, err := range errs {
		r = r.WithError(doctor.NewError("config", "[%s]: %s", err.Code, err.Message))
	}
	return r
}

func solveProject(_ context.Context, r *doctor.Result, _ util.Logger) *doctor.Result {
	if r.Errors.Find("missing") != nil {
		r.AddSolution("run [projectforge create] in this directory")
	}
	if r.Errors.Find("invalid") != nil {
		r.AddSolution("the project file isn't valid JSON, not sure what you can do")
	}
	if r.Errors.Find("config") != nil {
		r.AddSolution(fmt.Sprintf("use the %s UI to configure your project", util.AppName))
	}
	return r
}
