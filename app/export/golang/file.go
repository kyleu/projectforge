package golang

import (
	"fmt"
	"strings"

	"github.com/kyleu/projectforge/app/file"
	"github.com/kyleu/projectforge/app/lib/filesystem"
)

type File struct {
	Package string   `json:"package,omitempty"`
	Path    []string `json:"path,omitempty"`
	Name    string   `json:"name"`
	Imports Imports  `json:"imports"`
	Blocks  Blocks   `json:"blocks"`
}

func NewFile(pkg string, path []string, fn string) *File {
	return &File{Package: pkg, Path: path, Name: fn}
}

func (f *File) AddImport(t ImportType, v string) {
	f.Imports = f.Imports.Add(&Import{Type: t, Value: v})
}

func (f *File) AddBlocks(b ...*Block) {
	f.Blocks = append(f.Blocks, b...)
}

func (f *File) Render() *file.File {
	var content []string
	add := func(s string, args ...interface{}) {
		content = append(content, fmt.Sprintf(s+"\n", args...))
	}

	add("package %s", f.Package)

	if len(f.Imports) > 0 {
		add(f.Imports.Render())
	}

	for _, b := range f.Blocks {
		add(b.Render())
	}

	n := f.Name
	if !strings.HasSuffix(f.Name, ".go") {
		n += ".go"
	}
	return &file.File{Path: f.Path, Name: n, Mode: filesystem.DefaultMode, Content: strings.Join(content, "\n")}
}
