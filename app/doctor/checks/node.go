package checks

import (
	"github.com/kyleu/projectforge/app/doctor"
)

var node = &doctor.Check{
	Key:     "node",
	Section: "build",
	Title:   "Node.js",
	Summary: "Builds our web assets",
	URL:     "https://nodejs.org",
	UsedBy:  "Build of [client] TypeScript project and css pipeline",
	Fn:      doctor.SimpleOut(".", "node", []string{"-v"}, checkNode),
	Solve:   solveNode,
}

func checkNode(r *doctor.Result, out string) *doctor.Result {
	return r
}

func solveNode(r *doctor.Result) *doctor.Result {
	if r.Errors.Find("missing") != nil || r.Errors.Find("exitcode") != nil {
		r.Solution = "Install [Node.js] using your platform's package manager"
	}
	return r
}
