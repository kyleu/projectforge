package checks

import (
	"context"

	"go.uber.org/zap"
	"projectforge.dev/projectforge/app/doctor"
)

var qtc = &doctor.Check{
	Key:     "qtc",
	Section: "build",
	Title:   "Quicktemplate",
	Summary: "Compiles HTML and SQL templates at build time",
	URL:     "https://github.com/valyala/quicktemplate/qtc",
	UsedBy:  "Main server build",
	Fn:      simpleOut(".", "qtc", []string{"--help"}, noop),
	Solve:   solveQTC,
}

func solveQTC(ctx context.Context, r *doctor.Result, logger *zap.SugaredLogger) *doctor.Result {
	if r.Errors.Find("missing") != nil || r.Errors.Find("exitcode") != nil {
		r.AddSolution("go get -u github.com/valyala/quicktemplate/qtc")
	}
	return r
}
