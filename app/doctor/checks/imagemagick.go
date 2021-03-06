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

func checkImageMagick(ctx context.Context, r *doctor.Result, out string) *doctor.Result {
	if !strings.Contains(out, "ImageMagick") {
		return r.WithError(doctor.NewError("invalid", "[convert] is not provided by ImageMagick"))
	}
	return r
}

func solveImageMagick(ctx context.Context, r *doctor.Result, logger util.Logger) *doctor.Result {
	if r.Errors.Find("missing") != nil || r.Errors.Find("exitcode") != nil {
		r.AddSolution("Install [imagemagick] using your platform's package manager")
	}
	return r
}
