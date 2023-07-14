package checks

import (
	"context"

	"projectforge.dev/projectforge/app/doctor"
	"projectforge.dev/projectforge/app/util"
)

var Node = &doctor.Check{
	Key:     "node",
	Section: "build",
	Title:   "Node.js",
	Summary: "Builds our web assets",
	URL:     "https://nodejs.org",
	UsedBy:  "Build of [client] TypeScript project and css pipeline",
	Core:    true,
	Fn:      simpleOut(".", "node", []string{"-v"}, noop),
	Solve:   solveNode,
}

func solveNode(_ context.Context, r *doctor.Result, _ util.Logger) *doctor.Result {
	if r.Errors.Find("missing") != nil || r.Errors.Find("exitcode") != nil {
		r.AddPackageSolution("Node.js", "nodejs")
	}
	return r
}
