package checks

import (
	"context"

	"projectforge.dev/projectforge/app/doctor"
	"projectforge.dev/projectforge/app/util"
)

var Git = &doctor.Check{
	Key:     "git",
	Section: "build",
	Title:   "Git",
	Summary: "It's source control, it should really be installed",
	URL:     "https://git-scm.com",
	UsedBy:  "[bin/build/release.sh]",
	Fn:      simpleOut(".", "git", []string{"version"}, noop),
	Solve:   solveGit,
}

func solveGit(_ context.Context, r *doctor.Result, _ util.Logger) *doctor.Result {
	if r.Errors.HasMissing() || r.Errors.HasExitCode() != nil {
		r.AddPackageSolution("git", "git")
	}
	return r
}
