package export

import (
	"fmt"
	"strings"

	"github.com/kyleu/projectforge/app/file"
	"github.com/kyleu/projectforge/app/filesystem"
)

type GoTemplate struct {
	Package string    `json:"package,omitempty"`
	Path    []string  `json:"path,omitempty"`
	Name    string    `json:"name"`
	Imports GoImports `json:"imports"`
	Blocks  Blocks    `json:"blocks"`
}

func NewGoTemplate(pkg string, path []string, fn string) *GoTemplate {
	return &GoTemplate{Package: pkg, Path: path, Name: fn}
}

func (f *GoTemplate) AddImport(t GoImportType, v string) {
	f.Imports = append(f.Imports, &GoImport{Type: t, Value: v})
}

func (f *GoTemplate) AddBlocks(b ...*Block) {
	f.Blocks = append(f.Blocks, b...)
}

func (f *GoTemplate) Render() *file.File {
	var content []string
	add := func(s string, args ...interface{}) {
		content = append(content, fmt.Sprintf(s+"\n", args...))
	}

	if len(f.Imports) > 0 {
		add(f.Imports.RenderHTML())
	}

	for _, b := range f.Blocks {
		add(b.Render())
	}

	n := f.Name
	return &file.File{Path: f.Path, Name: n, Mode: filesystem.DefaultMode, Content: strings.Join(content, "\n")}
}
