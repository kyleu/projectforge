package golang

import (
	"fmt"
	"strings"

	"projectforge.dev/projectforge/app/file"
	"projectforge.dev/projectforge/app/lib/filesystem"
)

type Template struct {
	Path    []string `json:"path,omitempty"`
	Name    string   `json:"name"`
	Imports Imports  `json:"imports"`
	Blocks  Blocks   `json:"blocks"`
}

func NewGoTemplate(path []string, fn string) *Template {
	return &Template{Path: path, Name: fn}
}

func (f *Template) AddImport(i ...*Import) {
	for _, imp := range i {
		var hit bool
		for _, x := range f.Imports {
			if x.Equals(imp) {
				hit = true
			}
		}
		if !hit {
			f.Imports = append(f.Imports, imp)
		}
	}
}

func (f *Template) AddBlocks(b ...*Block) {
	f.Blocks = append(f.Blocks, b...)
}

func (f *Template) Render(addHeader bool) (*file.File, error) {
	var content []string
	add := func(s string, args ...any) {
		content = append(content, fmt.Sprintf(s+"\n", args...))
	}

	if addHeader {
		if strings.HasSuffix(f.Name, ".sql") {
			content = append(content, fmt.Sprintf("-- %s", file.HeaderContent))
		} else {
			content = append(content, fmt.Sprintf("<!-- %s -->", file.HeaderContent))
		}
	}
	if len(f.Imports) > 0 {
		add(f.Imports.RenderHTML())
	}

	for _, b := range f.Blocks {
		add(b.Render())
	}

	n := f.Name
	return &file.File{Path: f.Path, Name: n, Mode: filesystem.DefaultMode, Content: strings.Join(content, "\n")}, nil
}
