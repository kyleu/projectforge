package csharp

import (
	"fmt"
	"strings"

	"projectforge.dev/projectforge/app/file"
	"projectforge.dev/projectforge/app/lib/filesystem"
	"projectforge.dev/projectforge/app/util"
)

type File struct {
	Namespace string   `json:"namespace,omitempty"`
	Path      []string `json:"path,omitempty"`
	Name      string   `json:"name"`
	Imports   Imports  `json:"imports"`
	Blocks    Blocks   `json:"blocks"`
}

func NewFile(ns string, path []string, fn string) *File {
	return &File{Namespace: ns, Path: path, Name: fn}
}

func (f *File) AddImport(i ...string) {
	f.Imports = append(f.Imports, i...)
}

func (f *File) AddBlocks(b ...*Block) {
	f.Blocks = append(f.Blocks, b...)
}

func (f *File) Render(addHeader bool) (*file.File, error) {
	linebreak := "\n"
	var content []string
	if addHeader {
		if f.Namespace == "" {
			content = append(content, fmt.Sprintf("// %s", file.HeaderContent))
		} else {
			content = append(content, fmt.Sprintf("// Namespace %s - %s", f.Namespace, file.HeaderContent))
		}
	}

	if len(f.Imports) > 0 {
		content = append(content, f.Imports.Render(linebreak))
	}

	for _, b := range f.Blocks {
		x, err := b.Render(linebreak)
		if err != nil {
			return nil, err
		}
		content = append(content, x, "")
	}

	n := f.Name
	if !strings.HasSuffix(f.Name, util.ExtCS) {
		n += util.ExtCS
	}
	return &file.File{Path: f.Path, Name: n, Mode: filesystem.DefaultMode, Content: strings.Join(content, linebreak)}, nil
}
