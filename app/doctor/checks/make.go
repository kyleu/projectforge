package checks

import (
	"github.com/kyleu/projectforge/app/doctor"
	"go.uber.org/zap"
)

var mke = &doctor.Check{
	Key:     "make",
	Section: "build",
	Title:   "Make",
	Summary: "Compiles the project",
	URL:     "https://www.gnu.org/software/make",
	UsedBy:  "Main server build",
	Fn:      doctor.SimpleOut(".", "make", []string{"--version"}, noop),
	Solve:   solveMake,
}

func solveMake(r *doctor.Result, logger *zap.SugaredLogger) *doctor.Result {
	if r.Errors.Find("missing") != nil || r.Errors.Find("exitcode") != nil {
		r.Solution = "You should really have make installed"
	}
	return r
}
