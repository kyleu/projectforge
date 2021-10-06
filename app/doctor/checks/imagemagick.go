package checks

import (
	"strings"

	"github.com/kyleu/projectforge/app/doctor"
)

var imagemagick = &doctor.Check{
	Key:     "imagemagick",
	Title:   "ImageMagick",
	Summary: "Renders SVGs for the icon pipeline",
	Fn:      doctor.SimpleOut(".", "convert", []string{"-version"}, checkImageMagick),
	Solve:   solveImageMagick,
}

func checkImageMagick(r *doctor.Result, out string) *doctor.Result {
	if !strings.Contains(out, "ImageMagick") {
		return r.WithError(doctor.NewError("invalid", "[convert] is not provided by ImageMagick"))
	}
	return r
}

func solveImageMagick(r *doctor.Result) (*doctor.Result, error) {
	return r, nil
}
