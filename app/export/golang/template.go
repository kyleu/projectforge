package golang

import (
	"fmt"
	"strings"

	"github.com/kyleu/projectforge/app/file"
	"github.com/kyleu/projectforge/app/lib/filesystem"
)

type Template struct {
	Package string   `json:"package,omitempty"`
	Path    []string `json:"path,omitempty"`
	Name    string   `json:"name"`
	Imports Imports  `json:"imports"`
	Blocks  Blocks   `json:"blocks"`
}

func NewGoTemplate(pkg string, path []string, fn string) *Template {
	return &Template{Package: pkg, Path: path, Name: fn}
}

func (f *Template) AddImport(t ImportType, v string) {
	f.Imports = append(f.Imports, &Import{Type: t, Value: v})
}

func (f *Template) AddBlocks(b ...*Block) {
	f.Blocks = append(f.Blocks, b...)
}

func (f *Template) Render() *file.File {
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
