package export

import (
	"fmt"
	"strings"

	"github.com/kyleu/projectforge/app/file"
	"github.com/kyleu/projectforge/app/filesystem"
)

type GoFile struct {
	Package string    `json:"package,omitempty"`
	Path    []string  `json:"path,omitempty"`
	Name    string    `json:"name"`
	Imports GoImports `json:"imports"`
	Blocks  Blocks    `json:"blocks"`
}

func NewGoFile(pkg string, path []string, fn string) *GoFile {
	return &GoFile{Package: pkg, Path: path, Name: fn}
}

func (f *GoFile) AddImport(t GoImportType, v string) {
	f.Imports = append(f.Imports, &GoImport{Type: t, Value: v})
}

func (f *GoFile) AddBlocks(b ...*Block) {
	f.Blocks = append(f.Blocks, b...)
}

func (f *GoFile) Render() *file.File {
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
