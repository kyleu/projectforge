package file

import (
	"projectforge.dev/projectforge/app/lib/filesystem"
	"projectforge.dev/projectforge/app/util"
)

type File struct {
	Type    Type                `json:"type"`
	Path    []string            `json:"path,omitempty"`
	Name    string              `json:"name"`
	Mode    filesystem.FileMode `json:"mode,omitzero"`
	Content string              `json:"-"`
}

func NewFile(path string, mode filesystem.FileMode, b []byte) *File {
	p, n := util.StringSplitPath(path)
	return &File{Type: getType(n), Path: util.StringSplitPathAndTrim(p), Name: n, Mode: mode, Content: string(b)}
}

func (f *File) FullPath() string {
	return util.StringFilePath(util.StringFilePath(f.Path...), f.Name)
}

func (f *File) Ext() string {
	return util.StringSplitLastOnly(f.Name, '.', true)
}

const (
	prefix = "$PF_"
	suffix = "$"
)
