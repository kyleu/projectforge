package checks

import (
	"go.uber.org/zap"

	"projectforge.dev/projectforge/app/doctor"
)

var git = &doctor.Check{
	Key:     "git",
	Section: "build",
	Title:   "Git",
	Summary: "It's source control, it should really be installed",
	URL:     "https://git-scm.com",
	UsedBy:  "[bin/build/release.sh]",
	Fn:      simpleOut(".", "git", []string{"version"}, noop),
	Solve:   solveGit,
}

func solveGit(r *doctor.Result, logger *zap.SugaredLogger) *doctor.Result {
	if r.Errors.Find("missing") != nil || r.Errors.Find("exitcode") != nil {
		r.AddSolution("https://git-scm.com")
	}
	return r
}
