package assets

import (
	"embed"
	"fmt"
	"mime"
	"path/filepath"
	"sync"

	"github.com/pkg/errors"

	"projectforge.dev/projectforge/app/util"
)

//go:embed *
var FS embed.FS

type Entry struct {
	Bytes []byte
	Mime  string
	Hash  string
}

var (
	cache   = map[string]*Entry{}
	cacheMu sync.Mutex
)

func Embed(path string) (*Entry, error) {
	if path == "embed.go" {
		return nil, errors.New("invalid asset")
	}
	if x, ok := cache[path]; ok {
		return x, nil
	}
	cacheMu.Lock()
	defer cacheMu.Unlock()
	data, err := FS.ReadFile(path)
	if err != nil {
		return nil, errors.Wrapf(err, "error reading asset at [%s]", path)
	}
	mt := mime.TypeByExtension(filepath.Ext(path))
	h := util.HashFNV128UUID(string(data))
	e := &Entry{Bytes: data, Mime: mt, Hash: h.String()[:8]}
	cache[path] = e
	return e, nil
}

func URL(path string, logger util.Logger) string {
	e, err := Embed(path)
	if err != nil {
		logger.Warnf("unable to find asset [%s]: %s", path, err.Error())
		return fmt.Sprintf("/assets/%s", path)
	}
	return fmt.Sprintf("/assets/%s?hash=%s", path, e.Hash)
}

func ScriptElement(path string, deferFlag bool, logger util.Logger) string {
	if deferFlag {
		return fmt.Sprintf("<script src=%q defer=\"defer\"></script>", URL(path, logger))
	}
	return fmt.Sprintf("<script src=%q></script>", URL(path, logger))
}

func StylesheetElement(path string, logger util.Logger) string {
	return fmt.Sprintf(`<link rel="stylesheet" media="screen" href=%q>`, URL(path, logger))
}
