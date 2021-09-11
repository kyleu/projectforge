package checks

import (
	"fmt"

	"github.com/kyleu/projectforge/app/doctor"
	"github.com/kyleu/projectforge/app/util"
)

var rsvg = &doctor.Check{
	Key:     "rsvg",
	Title:   "rsvg",
	Summary: "Renders SVGs to PNG for the icon pipeline",
	Fn: func(r *doctor.Result) *doctor.Result {
		exitCode, out, err := util.RunProcessSimple("rsvg-convert -v", ".")
		if err != nil {
			msg := "[rsvg] is not present on your computer"
			return r.WithError(doctor.NewError("missing", msg, "rsvg"))
		}
		if exitCode != 0 {
			msg := fmt.Sprintf("[rsvg] returned [%d] as an exit code", exitCode)
			return r.WithError(doctor.NewError("exitcode", msg, "rsvg", fmt.Sprint(exitCode), out))
		}
		return r
	},
}
