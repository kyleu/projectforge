package help

import (
	"embed"
	"io/fs"
	"strings"

	"github.com/pkg/errors"

	"{{{ .Package }}}/app/util"
)

//go:embed *.md
var FS embed.FS

func List() ([]string, error) {
	files, err := fs.ReadDir(FS, ".")
	if err != nil {
		return nil, err
	}
	ret := util.NewStringSlice(make([]string, 0, len(files)))
	for _, f := range files {
		if strings.HasSuffix(f.Name(), util.ExtMarkdown) {
			ret.Push(strings.TrimSuffix(f.Name(), util.ExtMarkdown))
		}
	}
	return ret.Slice, nil
}

func Content(path string) (string, error) {
	data, err := FS.ReadFile(path)
	if err != nil {
		return "", errors.Wrapf(err, "error reading doc asset at [%s]", path)
	}
	body := strings.TrimSpace(string(data))
	return body, nil
}
