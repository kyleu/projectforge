package file

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/kyleu/projectforge/app/util"
)

type File struct {
	Type     string      `json:"type"`
	Path     []string    `json:"path,omitempty"`
	Name     string      `json:"name"`
	Mode     os.FileMode `json:"mode,omitempty"`
	Content  string      `json:"-"`
	fullPath string
}

func NewFile(path string, mode os.FileMode, b []byte) *File {
	p, n := util.SplitStringLast(path, '/', true)
	if n == "" {
		n = p
		p = ""
	}
	return &File{Type: getType(n), Path: util.SplitAndTrim(p, "/"), Name: n, Mode: mode, Content: string(b)}
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

func (f *File) Apply(cs *Changeset) *File {
	if cs != nil {
		r := cs.Replacements.ToReplacer(prefix, suffix)
		f.Content = r.Replace(f.Content)
	}
	return &File{Type: f.Type, Path: f.Path, Name: f.Name, Mode: f.Mode, Content: f.Content}
}

func getType(n string) string {
	switch {
	case n == "Makefile":
		return "makefile"
	case strings.HasSuffix(n, ".go"):
		return "go"
	default:
		return "unknown type (" + n + ")"
	}
}
