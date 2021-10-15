package checks

import (
	"github.com/kyleu/projectforge/app/doctor"
	"github.com/kyleu/projectforge/app/util"
	"go.uber.org/zap"
)

var pf = &doctor.Check{
	Key:     util.AppKey,
	Section: "app",
	Title:   util.AppName,
	Summary: "Confirms that [" + util.AppKey + "] is available on the path",
	URL:     util.AppURL,
	UsedBy:  util.AppName,
	Fn:      doctor.SimpleOut(".", "projectforge", []string{"--help"}, noop),
	Solve:   solvePF,
}

func solvePF(r *doctor.Result, logger *zap.SugaredLogger) *doctor.Result {
	if r.Errors.Find("missing") != nil || r.Errors.Find("exitcode") != nil {
		r.Solution = "Install [" + util.AppName + "] by following the instructions at [" + util.AppURL + "]"
	}
	return r
}
