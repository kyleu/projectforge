package checks

import (
	"context"

	"projectforge.dev/projectforge/app/doctor"
	"projectforge.dev/projectforge/app/util"
)

var Homebrew = &doctor.Check{
	Key:       "homebrew",
	Section:   "build",
	Title:     "Homebrew",
	Summary:   "Used to install other dependencies",
	URL:       "https://brew.sh",
	UsedBy:    "Package manager for macOS",
	Platforms: []string{"darwin"},
	Core:      true,
	Fn:        simpleOut(".", "brew", []string{"--help"}, noop),
	Solve:     solveHomebrew,
}

func solveHomebrew(_ context.Context, r *doctor.Result, _ util.Logger) *doctor.Result {
	if r.Errors.Find("missing") != nil || r.Errors.Find("exitcode") != nil {
		r.AddSolution(`!/bin/bash -c "$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/HEAD/install.sh)"`)
	}
	return r
}
