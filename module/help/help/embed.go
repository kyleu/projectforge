package help

import (
	"embed"
	"io/fs"
	"strings"

	"github.com/gomarkdown/markdown"
	"github.com/pkg/errors"
)

//go:embed *.md
var FS embed.FS

func List() ([]string, error) {
	files, err := fs.ReadDir(FS, ".")
	if err != nil {
		return nil, err
	}
	ret := make([]string, 0, len(files))
	for _, f := range files {
		if strings.HasSuffix(f.Name(), ".md") {
			ret = append(ret, strings.TrimSuffix(f.Name(), ".md"))
		}
	}
	return ret, nil
}

func Content(path string) ([]byte, error) {
	data, err := FS.ReadFile(path)
	if err != nil {
		return nil, errors.Wrapf(err, "error reading doc asset at [%s]", path)
	}
	return data, nil
}

func HTML(path string) (string, string, error) {
	if !strings.HasSuffix(path, ".md") {
		path += ".md"
	}
	data, err := Content(path)
	if err != nil {
		return "", "", err
	}
	return string(data), string(markdown.ToHTML(data, nil, nil)), nil
}
