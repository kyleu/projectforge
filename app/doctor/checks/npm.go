package checks

import (
	"context"

	"projectforge.dev/projectforge/app/doctor"
	"projectforge.dev/projectforge/app/util"
)

var NPM = &doctor.Check{
	Key:     "node",
	Section: "build",
	Title:   "NPM.js",
	Summary: "Download JavaScript dependencies",
	URL:     "https://www.npmjs.com/",
	UsedBy:  "Build of [client] TypeScript project and css pipeline",
	Core:    true,
	Fn:      simpleOut(".", "npm", []string{"-v"}, noop),
	Solve:   solveNPM,
}

func solveNPM(_ context.Context, r *doctor.Result, _ util.Logger) *doctor.Result {
	if r.Errors.Find("missing") != nil || r.Errors.Find("exitcode") != nil {
		r.AddPackageSolution("NPM", "npm")
	}
	return r
}
