package file

import (
	"os"
	"path/filepath"

	"github.com/kyleu/projectforge/app/util"
	"go.uber.org/zap"
)

type File struct {
	Path     []string    `json:"path,omitempty"`
	Name     string      `json:"name"`
	Mode     os.FileMode `json:"mode,omitempty"`
	Content  string      `json:"-"`
	fullPath string
}

func NewFile(path string, mode os.FileMode, b []byte, logger *zap.SugaredLogger) *File {
	p, n := util.SplitStringLast(path, '/', true)
	if n == "" {
		n = p
		p = ""
	}
	c := string(b)
	return &File{Path: util.SplitAndTrim(p, "/"), Name: n, Mode: mode, Content: c}
}

func (f *File) FullPath() string {
	if f.fullPath == "" {
		f.fullPath = filepath.Join(append(f.Path, f.Name)...)
	}
	return f.fullPath
}

const (
	prefix = "$PF_"
	suffix = "$"
)
