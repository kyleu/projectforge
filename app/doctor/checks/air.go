package checks

import (
	"github.com/kyleu/projectforge/app/doctor"
	"go.uber.org/zap"
)

var air = &doctor.Check{
	Key:     "air",
	Section: "build",
	Title:   "Air",
	Summary: "Used to recompile the project when files change",
	URL:     "https://github.com/cosmtrek/air",
	UsedBy:  "[bin/dev.sh]",
	Fn:      doctor.SimpleOut(".", "air", []string{"--help"}, noop),
	Solve:   solveAir,
}

func solveAir(r *doctor.Result, logger *zap.SugaredLogger) *doctor.Result {
	if r.Errors.Find("missing") != nil || r.Errors.Find("exitcode") != nil {
		r.Solution = "go get -u github.com/cosmtrek/air"
	}
	return r
}
