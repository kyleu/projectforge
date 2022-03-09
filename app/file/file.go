package file

import (
	"os"
	"path/filepath"

	"go.uber.org/zap"
	"projectforge.dev/projectforge/app/util"
)

type File struct {
	Type    Type        `json:"type"`
	Path    []string    `json:"path,omitempty"`
	Name    string      `json:"name"`
	Mode    os.FileMode `json:"mode,omitempty"`
	Content string      `json:"-"`
}

func NewFile(path string, mode os.FileMode, b []byte, addHeader bool, logger *zap.SugaredLogger) *File {
	p, n := util.StringSplitLast(path, '/', true)
	if n == "" {
		n = p
		p = ""
	}
	t := getType(n)
	c := string(b)
	if addHeader {
		c = contentWithHeader(t, c, logger)
	}
	return &File{Type: t, Path: util.StringSplitAndTrim(p, "/"), Name: n, Mode: mode, Content: c}
}

func (f *File) FullPath() string {
	return filepath.Join(f.Path...) + "/" + f.Name
}

const (
	prefix = "$PF_"
	suffix = "$"
)
