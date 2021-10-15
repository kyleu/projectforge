package checks

import (
	"strings"

	"github.com/kyleu/projectforge/app/doctor"
	"go.uber.org/zap"
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

func checkImageMagick(r *doctor.Result, out string) *doctor.Result {
	if !strings.Contains(out, "ImageMagick") {
		return r.WithError(doctor.NewError("invalid", "[convert] is not provided by ImageMagick"))
	}
	return r
}

func solveImageMagick(r *doctor.Result, logger *zap.SugaredLogger) *doctor.Result {
	if r.Errors.Find("missing") != nil || r.Errors.Find("exitcode") != nil {
		r.AddSolution("Install [imagemagick] using your platform's package manager")
	}
	return r
}
