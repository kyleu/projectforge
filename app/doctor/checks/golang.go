package checks

import (
	"context"
	"strings"

	"golang.org/x/mod/semver"

	"projectforge.dev/projectforge/app/doctor"
	"projectforge.dev/projectforge/app/util"
)

const golangVersion = "v1.20.0"

var Golang = &doctor.Check{
	Key:     "golang",
	Section: "build",
	Title:   "Go",
	Summary: "The main programming language",
	URL:     "https://golang.org",
	UsedBy:  "All builds",
	Core:    true,
	Fn: simpleOut(".", "go", []string{"version"}, func(_ context.Context, r *doctor.Result, out string) *doctor.Result {
		if r.Status == util.KeyError {
			return r
		}
		startIdx := strings.Index(out, "go1.")
		if startIdx == -1 {
			return r.WithError(&doctor.Error{Code: util.KeyUnknown, Message: "can't parse result of [go version]"})
		}
		endIdx := strings.LastIndex(out, " ")
		if endIdx == -1 || endIdx <= startIdx {
			return r.WithError(&doctor.Error{Code: util.KeyUnknown, Message: "can't parse end result of [go version]"})
		}
		v := "v" + out[startIdx+2:endIdx]
		if semver.Compare(v, golangVersion) < 0 {
			return r.WithError(&doctor.Error{Code: "minversion", Message: "Go version [" + v + "] must be equal or higher than [" + golangVersion + "]"})
		}
		return r
	}),
	Solve: solveGo,
}

func solveGo(_ context.Context, r *doctor.Result, _ util.Logger) *doctor.Result {
	if r.Errors.Find("missing") != nil || r.Errors.Find("exitcode") != nil || r.Errors.Find("minversion") != nil {
		r.AddPackageSolution("Go", "golang")
	}
	return r
}
