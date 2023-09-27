// Package doc - Content managed by Project Forge, see [projectforge.md] for details.
package doc

import (
	"embed"

	"github.com/gomarkdown/markdown"
	"github.com/pkg/errors"
)

//go:embed *
var FS embed.FS

func Content(path string) ([]byte, error) {
	if path == "embed.go" {
		return nil, errors.New("invalid asset")
	}
	data, err := FS.ReadFile(path)
	if err != nil {
		return nil, errors.Wrapf(err, "error reading doc asset at [%s]", path)
	}
	return data, nil
}

var (
	titleCache = map[string]string{}
	htmlCache  = map[string]string{}
)

func HTML(key string, path string, f func(s string) (string, string, error)) (string, string, error) {
	if curr, ok := htmlCache[key]; ok {
		return titleCache[key], curr, nil
	}
	data, err := Content(path)
	if err != nil {
		return "", "", err
	}
	return HTMLString(key, data, f)
}

func HTMLString(key string, data []byte, f func(s string) (string, string, error)) (string, string, error) {
	if curr, ok := htmlCache[key]; ok {
		return titleCache[key], curr, nil
	}
	var title string
	html := string(markdown.ToHTML(data, nil, nil))
	if f != nil {
		var err error
		title, html, err = f(html)
		if err != nil {
			return "", "", err
		}
	}
	titleCache[key] = title
	htmlCache[key] = html
	return title, html, nil
}
