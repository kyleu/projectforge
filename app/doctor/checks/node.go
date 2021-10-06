package checks

import (
	"github.com/kyleu/projectforge/app/doctor"
)

var node = &doctor.Check{
	Key:     "node",
	Title:   "NodeJS",
	Summary: "Builds our web assets",
	Fn:      doctor.SimpleOut(".", "node", []string{"-v"}, checkNode),
	Solve:   solveNode,
}

func checkNode(r *doctor.Result, out string) *doctor.Result {
	return r
}

func solveNode(r *doctor.Result) (*doctor.Result, error) {
	return r, nil
}
