package checks

import (
	"context"

	"projectforge.dev/projectforge/app/doctor"
	"projectforge.dev/projectforge/app/util"
)

var Golangcilint = &doctor.Check{
	Key:     "golangcilint",
	Section: "build",
	Title:   "golangci-lint",
	Summary: "Check for style and linting errors",
	URL:     "https://golangci-lint.run",
	UsedBy:  "[bin/check.sh]",
	Fn:      simpleOut(".", "golangci-lint", nil, noop),
	Solve:   solveGolangcilint,
}

func solveGolangcilint(_ context.Context, r *doctor.Result, _ util.Logger) *doctor.Result {
	if r.Errors.HasMissing() || r.Errors.HasExitCode() {
		r.AddSolution("!go install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@latest")
	}
	return r
}
