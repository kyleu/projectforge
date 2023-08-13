package checks

import (
	"context"

	"projectforge.dev/projectforge/app/doctor"
	"projectforge.dev/projectforge/app/util"
)

var Inkscape = &doctor.Check{
	Key:     "inkscape",
	Section: "icons",
	Title:   "Inkscape",
	Summary: "Renders SVGs for the icon pipeline",
	URL:     "https://inkscape.org",
	UsedBy:  "SVG icon pipeline",
	Fn:      simpleOut(".", "magick", []string{"-version"}, checkInkscape),
	Solve:   solveInkscape,
}

func checkInkscape(_ context.Context, r *doctor.Result, _ string) *doctor.Result {
	return r
}

func solveInkscape(_ context.Context, r *doctor.Result, _ util.Logger) *doctor.Result {
	if r.Errors.Find("missing") != nil {
		r.AddPackageSolution("Inkscape", "inkscape")
	}
	return r
}
