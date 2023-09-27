package golang

import (
	"fmt"
	"strings"

	"github.com/samber/lo"

	"projectforge.dev/projectforge/app/file"
	"projectforge.dev/projectforge/app/lib/filesystem"
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

func (f *File) AddImport(i ...*Import) {
	f.Imports = f.Imports.Add(i...)
}

func (f *File) AddBlocks(b ...*Block) {
	f.Blocks = append(f.Blocks, b...)
}

func (f *File) Render(addHeader bool, linebreak string) (*file.File, error) {
	var content []string
	add := func(s string, args ...any) {
		content = append(content, fmt.Sprintf(s+linebreak, args...))
	}

	if addHeader {
		if f.Package == "" {
			content = append(content, fmt.Sprintf("// %s", file.HeaderContent))
		} else {
			content = append(content, fmt.Sprintf("// Package %s - %s", f.Package, file.HeaderContent))
		}
	}
	add("package %s", f.Package)

	if len(f.Imports) > 0 {
		add(f.Imports.Render(linebreak))
	}

	lo.ForEach(f.Blocks, func(b *Block, _ int) {
		add(b.Render(linebreak))
	})

	n := f.Name
	if !strings.HasSuffix(f.Name, ".go") {
		n += ".go"
	}
	return &file.File{Path: f.Path, Name: n, Mode: filesystem.DefaultMode, Content: strings.Join(content, linebreak)}, nil
}
