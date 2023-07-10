package checks

import (
	"context"
	"strings"

	"projectforge.dev/projectforge/app/doctor"
	"projectforge.dev/projectforge/app/util"
)

var Imagemagick = &doctor.Check{
	Key:     "imagemagick",
	Section: "icons",
	Title:   "ImageMagick",
	Summary: "Renders SVGs for the icon pipeline",
	URL:     "https://imagemagick.org",
	UsedBy:  "SVG icon pipeline",
	Fn:      simpleOut(".", "magick", []string{"-version"}, checkImageMagick),
	Solve:   solveImageMagick,
}

func checkImageMagick(_ context.Context, r *doctor.Result, out string) *doctor.Result {
	if !strings.Contains(out, "ImageMagick") {
		return r.WithError(doctor.NewError("invalid", "[magick] is not provided by ImageMagick"))
	}
	return r
}

func solveImageMagick(_ context.Context, r *doctor.Result, _ util.Logger) *doctor.Result {
	if r.Errors.Find("missing") != nil {
		r.AddPackageSolution("Imagemagick", "imagemagick")
	}
	return r
}
