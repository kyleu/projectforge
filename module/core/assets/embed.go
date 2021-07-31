package assets

import (
	"bytes"
	"compress/gzip"
	"crypto/md5" // nolint
	"embed"
	"encoding/hex"
	"mime"
	"path/filepath"

	"github.com/pkg/errors"
)

//go:embed *
var assetFS embed.FS

func EmbedAsset(path string) ([]byte, string, string, error) {
	var b bytes.Buffer

	data, err := assetFS.ReadFile(path)
	if err != nil {
		return nil, "", "", errors.Wrapf(err, "error reading asset at [%s]", path)
	}

	if data != nil {
		w := gzip.NewWriter(&b)
		_, _ = w.Write(data)
		_ = w.Close()
		data = b.Bytes()
	}

	// nolint
	sum := md5.Sum(data)

	return data, hex.EncodeToString(sum[1:]), mime.TypeByExtension(filepath.Ext(path)), nil
}
