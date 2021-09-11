package checks

import (
	"fmt"

	"github.com/kyleu/projectforge/app/doctor"
	"github.com/kyleu/projectforge/app/util"
)

var node = &doctor.Check{
	Key:     "node",
	Title:   "NodeJS",
	Summary: "Builds our web assets",
	Fn: func(r *doctor.Result) *doctor.Result {
		exitCode, out, err := util.RunProcessSimple("node -v", ".")
		if err != nil {
			msg := "[node] is not present on your computer"
			return r.WithError(doctor.NewError("missing", msg, "node"))
		}
		if exitCode != 0 {
			msg := fmt.Sprintf("[node] returned [%d] as an exit code", exitCode)
			return r.WithError(doctor.NewError("exitcode", msg, "node", fmt.Sprint(exitCode), out))
		}
		return r
	},
}
