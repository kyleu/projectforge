package checks

import (
	"github.com/kyleu/projectforge/app/doctor"
)

var qtc = &doctor.Check{
	Key:     "qtc",
	Section: "build",
	Title:   "Quicktemplate",
	Summary: "Compiles HTML and SQL templates at build time",
	URL:     "https://github.com/valyala/quicktemplate/qtc",
	UsedBy:  "Main server build",
	Fn:      doctor.SimpleOut(".", "qtc", []string{"--help"}, checkQTC),
	Solve:   solveQTC,
}

func checkQTC(r *doctor.Result, out string) *doctor.Result {
	return r
}

func solveQTC(r *doctor.Result) *doctor.Result {
	if r.Errors.Find("missing") != nil || r.Errors.Find("exitcode") != nil {
		r.Solution = "go get -u github.com/valyala/quicktemplate/qtc"
	}
	return r
}
