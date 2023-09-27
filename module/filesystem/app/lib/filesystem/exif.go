package filesystem

import (
	"strings"

	"github.com/dsoprea/go-exif/v3"
	"github.com/pkg/errors"
	"github.com/samber/lo"

	"{{{ .Package }}}/app/util"
)

func ExifExtract(b []byte) (util.ValueMap, error) {
	x, err := exif.SearchAndExtractExif(b)
	if err != nil {
		return nil, errors.Wrapf(err, "unable to extract exif data")
	}
	entries, _, err := exif.GetFlatExifDataUniversalSearch(x, nil, true)
	if err != nil {
		return nil, errors.Wrapf(err, "unable to extract exif data")
	}
	ret := lo.SliceToMap(entries, func(t exif.ExifTag) (string, any) {
		return t.TagName, strings.TrimSuffix(strings.TrimPrefix(t.Formatted, "["), "]")
	})
	return ret, nil
}

func ImageType(path ...string) string {
	if len(path) == 0 {
		return ""
	}
	ext := path[len(path)-1]
	if !strings.Contains(ext, ".") {
		return ""
	}
	t := ext[strings.LastIndex(ext, ".")+1:]
	switch t {
	case "bmp", "gif", "jpg", "jpeg", "png":
		return t
	default:
		return ""
	}
}
