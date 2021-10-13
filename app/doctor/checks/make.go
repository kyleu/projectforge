package checks

import (
	"github.com/kyleu/projectforge/app/doctor"
)

var mke = &doctor.Check{
	Key:     "make",
	Section: "build",
	Title:   "Make",
	Summary: "Compiles the project",
	URL:     "https://www.gnu.org/software/make",
	UsedBy:  "Main server build",
	Fn:      doctor.SimpleOut(".", "make", []string{"--version"}, checkMake),
	Solve:   solveMake,
}

func checkMake(r *doctor.Result, out string) *doctor.Result {
	return r
}

func solveMake(r *doctor.Result) *doctor.Result {
	if r.Errors.Find("missing") != nil || r.Errors.Find("exitcode") != nil {
		r.Solution = "You should really have make installed"
	}
	return r
}
