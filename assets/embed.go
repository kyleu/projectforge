// Content managed by Project Forge, see [projectforge.md] for details.
package assets

import (
	"embed"
	"mime"
	"path/filepath"

	"github.com/pkg/errors"
)

//go:embed *
var FS embed.FS

func EmbedAsset(path string) ([]byte, string, error) {
	if path == "embed.go" {
		return nil, "", errors.New("invalid asset")
	}
	data, err := FS.ReadFile(path)
	if err != nil {
		return nil, "", errors.Wrapf(err, "error reading asset at [%s]", path)
	}

	return data, mime.TypeByExtension(filepath.Ext(path)), nil
}
