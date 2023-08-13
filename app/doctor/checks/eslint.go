package checks

import (
	"context"

	"projectforge.dev/projectforge/app/doctor"
	"projectforge.dev/projectforge/app/util"
)

var ESLint = &doctor.Check{
	Key:     "eslint",
	Section: "icons",
	Title:   "ESLint",
	Summary: "Lints the TypeScript codebase",
	URL:     "https://eslint.org",
	UsedBy:  "bin/check-client.sh",
	Fn:      simpleOut(".", "eslint", []string{"--help"}, checkESLint),
	Solve:   solveESLint,
}

func checkESLint(_ context.Context, r *doctor.Result, _ string) *doctor.Result {
	return r
}

func solveESLint(_ context.Context, r *doctor.Result, _ util.Logger) *doctor.Result {
	if r.Errors.Find("missing") != nil {
		r.AddPackageSolution("ESLint", "eslint")
	}
	return r
}
