package checks

import (
	"context"
	"strings"

	"projectforge.dev/projectforge/app/doctor"
	"projectforge.dev/projectforge/app/util"
)

var Air = &doctor.Check{
	Key:       "air",
	Section:   "build",
	Title:     "Air",
	Summary:   "Used to recompile the project when files change",
	URL:       "https://github.com/cosmtrek/air",
	UsedBy:    "[bin/dev.sh]",
	Platforms: []string{"!windows"},
	Fn: simpleOut(".", "air", []string{"--help"}, func(ctx context.Context, r *doctor.Result, out string) *doctor.Result {
		if strings.Contains(out, "Command 'air' not found") {
			return r.WithError(doctor.NewError("missing", "[air] is not present on your computer"))
		}
		return r
	}),
	Solve: solveAir,
}

func solveAir(_ context.Context, r *doctor.Result, _ util.Logger) *doctor.Result {
	if r.Errors.Find("missing") != nil || r.Errors.Find("exitcode") != nil {
		r.AddSolution("!go install github.com/cosmtrek/air@latest")
	}
	return r
}
