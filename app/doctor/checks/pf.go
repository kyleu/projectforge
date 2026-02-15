package checks

import (
	"context"
	"fmt"

	"projectforge.dev/projectforge/app/doctor"
	"projectforge.dev/projectforge/app/util"
)

var PF = &doctor.Check{
	Key:     util.AppKey,
	Section: "app",
	Title:   util.AppName,
	Summary: "Confirms that [" + util.AppKey + "] is available on the path",
	URL:     util.AppURL,
	UsedBy:  util.AppName,
	Fn:      simpleOut(".", "projectforge", []string{"--help"}, noop),
	Solve:   solvePF,
}

func solvePF(_ context.Context, r *doctor.Result, _ util.Logger) *doctor.Result {
	if r.Errors.Find("missing") != nil || r.Errors.Find("exit-code") != nil {
		msg := "Install [%s] by following the instructions at [%s], and ensure [%s] is on your path"
		r.AddSolution(fmt.Sprintf(msg, util.AppName, util.AppURL, util.AppKey))
	}
	return r
}
