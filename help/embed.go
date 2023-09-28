// Package help - Content managed by Project Forge, see [projectforge.md] for details.
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

func Content(path string) (string, error) {
	data, err := FS.ReadFile(path)
	if err != nil {
		return "", errors.Wrapf(err, "error reading doc asset at [%s]", path)
	}
	body := strings.TrimSpace(string(data))
	return body, nil
}

func HTML(path string) (string, string, error) {
	if !strings.HasSuffix(path, ".md") {
		path += ".md"
	}
	data, err := Content(path)
	if err != nil {
		return "", "", err
	}
	return data, string(markdown.ToHTML([]byte(data), nil, nil)), nil
}
