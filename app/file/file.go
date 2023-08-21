package file

import (
	"path/filepath"

	"projectforge.dev/projectforge/app/lib/filesystem"
	"projectforge.dev/projectforge/app/util"
)

type File struct {
	Type    Type                `json:"type"`
	Path    []string            `json:"path,omitempty"`
	Name    string              `json:"name"`
	Mode    filesystem.FileMode `json:"mode,omitempty"`
	Content string              `json:"-"`
}

func NewFile(path string, mode filesystem.FileMode, b []byte, addHeader bool, logger util.Logger) *File {
	p, n := util.StringSplitPath(path)
	t := getType(n)
	c := string(b)
	if addHeader {
		c = contentWithHeader(path, t, c, util.StringDetectLinebreak(c), logger)
	}
	return &File{Type: t, Path: util.StringSplitPathAndTrim(p), Name: n, Mode: mode, Content: c}
}

func (f *File) FullPath() string {
	return filepath.Join(filepath.Join(f.Path...), f.Name)
}

func (f *File) Ext() string {
	return util.StringSplitLastOnly(f.Name, '.', true)
}

const (
	prefix = "$PF_"
	suffix = "$"
)
