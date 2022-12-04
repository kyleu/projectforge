// Content managed by Project Forge, see [projectforge.md] for details.
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

var htmlCache = map[string]string{}

func HTML(path string, f func(s string) (string, error)) (string, error) {
	if curr, ok := htmlCache[path]; ok {
		return curr, nil
	}
	data, err := Content(path)
	if err != nil {
		return "", err
	}
	return HTMLString(path, data, f)
}

func HTMLString(key string, data []byte, f func(s string) (string, error)) (string, error) {
	if curr, ok := htmlCache[key]; ok {
		return curr, nil
	}
	html := string(markdown.ToHTML(data, nil, nil))
	if f != nil {
		var err error
		html, err = f(html)
		if err != nil {
			return "", err
		}
	}
	htmlCache[key] = html
	return html, nil
}
