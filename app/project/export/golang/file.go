package golang

import (
	"fmt"
	"strings"

	"projectforge.dev/projectforge/app/file"
	"projectforge.dev/projectforge/app/lib/filesystem"
	"projectforge.dev/projectforge/app/lib/metamodel/model"
	"projectforge.dev/projectforge/app/util"
)

type File struct {
	Package string        `json:"package,omitempty"`
	Path    []string      `json:"path,omitempty"`
	Name    string        `json:"name"`
	Imports model.Imports `json:"imports"`
	Blocks  Blocks        `json:"blocks"`
}

func NewFile(pkg string, path []string, fn string) *File {
	return &File{Package: pkg, Path: path, Name: fn}
}

func (f *File) AddImport(i ...*model.Import) {
	f.Imports = f.Imports.Add(i...)
}

func (f *File) AddBlocks(b ...*Block) {
	f.Blocks = append(f.Blocks, b...)
}

func (f *File) Render(linebreak string) (*file.File, error) {
	var content []string
	add := func(s string, args ...any) {
		content = append(content, fmt.Sprintf(s+linebreak, args...))
	}

	add("package %s", f.Package)

	if len(f.Imports) > 0 {
		add(f.Imports.Render(linebreak))
	}

	for _, b := range f.Blocks {
		x, err := b.Render(linebreak)
		if err != nil {
			return nil, err
		}
		add(x)
	}

	n := f.Name
	if !strings.HasSuffix(f.Name, util.ExtGo) {
		n += util.ExtGo
	}
	return &file.File{Path: f.Path, Name: n, Mode: filesystem.DefaultMode, Content: util.StringJoin(content, linebreak)}, nil
}
