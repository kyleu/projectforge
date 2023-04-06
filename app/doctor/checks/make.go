package checks

import (
	"context"

	"projectforge.dev/projectforge/app/doctor"
	"projectforge.dev/projectforge/app/util"
)

var mke = &doctor.Check{
	Key:     "make",
	Section: "build",
	Title:   "Make",
	Summary: "Compiles the project",
	URL:     "https://www.gnu.org/software/make",
	UsedBy:  "Main server build",
	Fn:      simpleOut(".", "make", []string{"--version"}, noop),
	Solve:   solveMake,
}

func solveMake(_ context.Context, r *doctor.Result, _ util.Logger) *doctor.Result {
	if r.Errors.Find("missing") != nil || r.Errors.Find("exitcode") != nil {
		r.AddSolution("You should really have make installed")
	}
	return r
}
