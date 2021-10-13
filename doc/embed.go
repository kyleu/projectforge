package doc

import (
	"embed"

	"github.com/gomarkdown/markdown"
	"github.com/pkg/errors"
)

//go:embed *
var assetFS embed.FS

func Content(path string) ([]byte, error) {
	data, err := assetFS.ReadFile(path)
	if err != nil {
		return nil, errors.Wrapf(err, "error reading doc asset at [%s]", path)
	}

	return data, nil
}

var htmlCache = map[string]string{}

func HTML(path string) (string, error) {
	if curr, ok := htmlCache[path]; ok {
		return curr, nil
	}
	data, err := Content(path)
	if err != nil {
		return "", err
	}
	html := string(markdown.ToHTML(data, nil, nil))
	htmlCache[path] = html
	return html, nil
}
