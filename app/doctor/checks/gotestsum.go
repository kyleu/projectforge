package checks

import (
	"context"

	"projectforge.dev/projectforge/app/doctor"
	"projectforge.dev/projectforge/app/util"
)

var Gotestsum = &doctor.Check{
	Key:     "gotestsum",
	Section: "test",
	Title:   "gotestsum",
	Summary: "Runs the tests and displays the results",
	URL:     "https://github.com/gotestyourself/gotestsum",
	UsedBy:  "[bin/test.sh]",
	Fn:      simpleOut(".", "gotestsum", []string{"--version"}, noop),
	Solve:   solveGotestsum,
}

func solveGotestsum(_ context.Context, r *doctor.Result, _ util.Logger) *doctor.Result {
	if r.Errors.HasMissing() || r.Errors.HasExitCode() {
		r.AddSolution("!go install gotest.tools/gotestsum@latest")
	}
	return r
}
