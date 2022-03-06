package checks

import (
	"go.uber.org/zap"
	"projectforge.dev/app/doctor"
)

var node = &doctor.Check{
	Key:     "node",
	Section: "build",
	Title:   "Node.js",
	Summary: "Builds our web assets",
	URL:     "https://nodejs.org",
	UsedBy:  "Build of [client] TypeScript project and css pipeline",
	Fn:      simpleOut(".", "node", []string{"-v"}, noop),
	Solve:   solveNode,
}

func solveNode(r *doctor.Result, logger *zap.SugaredLogger) *doctor.Result {
	if r.Errors.Find("missing") != nil || r.Errors.Find("exitcode") != nil {
		r.AddSolution("Install [Node.js] using your platform's package manager")
	}
	return r
}
