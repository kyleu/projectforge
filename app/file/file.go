package file

import (
	"os"
	"path/filepath"

	"projectforge.dev/projectforge/app/util"
)

type File struct {
	Type    Type        `json:"type"`
	Path    []string    `json:"path,omitempty"`
	Name    string      `json:"name"`
	Mode    os.FileMode `json:"mode,omitempty"`
	Content string      `json:"-"`
}

func NewFile(path string, mode os.FileMode, b []byte, addHeader bool, logger util.Logger) *File {
	p, n := util.StringSplitLast(path, '/', true)
	if n == "" {
		n = p
		p = ""
	}
	t := getType(n)
	c := string(b)
	if addHeader {
		c = contentWithHeader(path, t, c, logger)
	}
	return &File{Type: t, Path: util.StringSplitAndTrim(p, "/"), Name: n, Mode: mode, Content: c}
}

func (f *File) FullPath() string {
	return filepath.Join(f.Path...) + "/" + f.Name
}

func (f *File) Ext() string {
	l, r := util.StringSplitLast(f.Name, '.', true)
	if r == "" {
		return l
	}
	return r
}

const (
	prefix = "$PF_"
	suffix = "$"
)
