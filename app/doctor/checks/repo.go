package checks

import (
	"context"
	"strings"

	"projectforge.dev/projectforge/app/doctor"
	"projectforge.dev/projectforge/app/lib/telemetry"
	"projectforge.dev/projectforge/app/util"
)

var Repo = &doctor.Check{
	Key:     "repo",
	Section: "release",
	Title:   "Git Repo",
	Summary: "Verifies the git repository in the working directory",
	URL:     "https://git-scm.com",
	UsedBy:  "bin/build/release.sh",
	Fn:      checkRepo,
	Solve:   solveRepo,
}

func checkRepo(ctx context.Context, r *doctor.Result, logger util.Logger) *doctor.Result {
	exitCode, _, err := telemetry.RunProcessSimple(ctx, "git status", ".", logger)
	if err != nil {
		return r.WithError(doctor.NewError("missing", "[git] is not present on your computer"))
	}
	if exitCode == 128 {
		return r.WithError(doctor.NewError("norepo", "no git repository in current directory"))
	}
	exitCode, out, err := telemetry.RunProcessSimple(ctx, "git status", ".", logger)
	if err != nil {
		return r.WithError(doctor.NewError(util.KeyError, "can't run [git status]: %+v", err))
	}
	if exitCode != 0 {
		return r.WithError(doctor.NewError(util.KeyError, "can't run [git status]: %s", out))
	}
	if out = strings.TrimSpace(out); out == "" {
		return r.WithError(doctor.NewError("noremote", "no git remote configured", out))
	}
	exitCode, _, err = telemetry.RunProcessSimple(ctx, "git log -1", ".", logger)
	if err != nil {
		return r.WithError(doctor.NewError(util.KeyError, "can't run [git log]: %+v", err))
	}
	if exitCode == 128 {
		return r.WithError(doctor.NewError("nocommit", "git repo must have at least one commit"))
	}

	return r
}

func solveRepo(_ context.Context, r *doctor.Result, _ util.Logger) *doctor.Result {
	if r.Errors.Find("norepo") != nil {
		r.AddSolution("!git init")
	}
	if r.Errors.Find("noremote") != nil {
		p, _, dr := loadRootProject(r)
		dr.AddSolution("!git remote add origin " + p.Info.Sourcecode + ".git")
	}
	if r.Errors.Find("nocommit") != nil {
		r.AddSolution("!git commit -am \"initial commit\"")
	}
	return r
}
