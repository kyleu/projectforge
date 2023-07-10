package checks

import (
	"context"

	"projectforge.dev/projectforge/app/doctor"
	"projectforge.dev/projectforge/app/util"
)

var Make = &doctor.Check{
	Key:       "make",
	Section:   "build",
	Title:     "Make",
	Summary:   "Compiles the project",
	URL:       "https://www.gnu.org/software/make",
	UsedBy:    "Main server build",
	Platforms: []string{"!windows"},
	Fn:        simpleOut(".", "make", []string{"--version"}, noop),
	Solve:     solveMake,
}

func solveMake(_ context.Context, r *doctor.Result, _ util.Logger) *doctor.Result {
	if r.Errors.Find("missing") != nil || r.Errors.Find("exitcode") != nil {
		r.AddPackageSolution("Make", "make")
	}
	return r
}
