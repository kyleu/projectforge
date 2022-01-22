// Content managed by Project Forge, see [projectforge.md] for details.
package assets

import (
	"embed"
	"mime"
	"path/filepath"

	"github.com/pkg/errors"
)

//go:embed *
var assetFS embed.FS

func EmbedAsset(path string) ([]byte, string, error) {
	data, err := assetFS.ReadFile(path)
	if err != nil {
		return nil, "", errors.Wrapf(err, "error reading asset at [%s]", path)
	}

	return data, mime.TypeByExtension(filepath.Ext(path)), nil
}
