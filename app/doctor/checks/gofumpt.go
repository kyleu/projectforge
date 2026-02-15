package checks

import (
	"context"

	"projectforge.dev/projectforge/app/doctor"
	"projectforge.dev/projectforge/app/util"
)

var Gofumpt = &doctor.Check{
	Key:     "gofumpt",
	Section: "build",
	Title:   "gofumpt",
	Summary: "Formats the code; stricter than [go fmt]",
	URL:     "https://github.com/mvdan/gofumpt",
	UsedBy:  "[bin/format.sh]",
	Fn:      simpleOut(".", "gofumpt", []string{"--version"}, noop),
	Solve:   solveGofumpt,
}

func solveGofumpt(_ context.Context, r *doctor.Result, _ util.Logger) *doctor.Result {
	if r.Errors.Find("missing") != nil || r.Errors.Find("exit-code") != nil {
		r.AddSolution("!go install mvdan.cc/gofumpt@latest")
	}
	return r
}
