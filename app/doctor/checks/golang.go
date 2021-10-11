package checks

import (
	"github.com/kyleu/projectforge/app/doctor"
)

var golang = &doctor.Check{
	Key:     "golang",
	Section: "build",
	Title:   "Go",
	Summary: "The main programming language",
	URL:     "https://golang.org",
	UsedBy:  "All builds",
	Fn:      doctor.SimpleOut(".", "go", []string{"version"}, checkGo),
	Solve:   solveGo,
}

func checkGo(r *doctor.Result, out string) *doctor.Result {
	return r
}

func solveGo(r *doctor.Result) *doctor.Result {
	if r.Errors.Find("missing") != nil || r.Errors.Find("exitcode") != nil {
		r.Solution = "Download Go for your platform"
	}
	return r
}
