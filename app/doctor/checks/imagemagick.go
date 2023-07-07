package checks

import (
	"context"
	"strings"

	"projectforge.dev/projectforge/app/doctor"
	"projectforge.dev/projectforge/app/util"
)

var imagemagick = &doctor.Check{
	Key:     "imagemagick",
	Section: "icons",
	Title:   "ImageMagick",
	Summary: "Renders SVGs for the icon pipeline",
	URL:     "https://imagemagick.org",
	UsedBy:  "SVG icon pipeline",
	Fn:      simpleOut(".", "convert", []string{"-version"}, checkImageMagick),
	Solve:   solveImageMagick,
}

func checkImageMagick(_ context.Context, r *doctor.Result, out string) *doctor.Result {
	if !strings.Contains(out, "ImageMagick") {
		return r.WithError(doctor.NewError("invalid", "[convert] is not provided by ImageMagick"))
	}
	return r
}

func solveImageMagick(_ context.Context, r *doctor.Result, _ util.Logger) *doctor.Result {
	if r.Errors.Find("exitcode") != nil {
		r.AddSolution("On Windows 10+, Imagemagick's [convert] program is superseded by a Windows utility, delete it or ignore the error")
	}
	if r.Errors.Find("missing") != nil {
		r.AddPackageSolution("Imagemagick", "imagemagick")
	}
	return r
}
