package golang

import (
	"fmt"

	"github.com/samber/lo"

	"projectforge.dev/projectforge/app/file"
	"projectforge.dev/projectforge/app/lib/filesystem"
	"projectforge.dev/projectforge/app/lib/metamodel/model"
	"projectforge.dev/projectforge/app/util"
)

type Template struct {
	Path    []string      `json:"path,omitempty"`
	Name    string        `json:"name"`
	Imports model.Imports `json:"imports"`
	Blocks  Blocks        `json:"blocks"`
}

func NewGoTemplate(path []string, fn string) *Template {
	return &Template{Path: path, Name: fn}
}

func (f *Template) AddImport(i ...*model.Import) {
	lo.ForEach(i, func(imp *model.Import, _ int) {
		hit := lo.ContainsBy(f.Imports, func(x *model.Import) bool {
			return x.Equals(imp)
		})
		if !hit {
			f.Imports = append(f.Imports, imp)
		}
	})
}

func (f *Template) AddBlocks(b ...*Block) {
	f.Blocks = append(f.Blocks, b...)
}

func (f *Template) Render(linebreak string) (*file.File, error) {
	var content []string
	add := func(s string, args ...any) {
		content = append(content, fmt.Sprintf(s+linebreak, args...))
	}

	if len(f.Imports) > 0 {
		add(f.Imports.RenderHTML(linebreak))
	}

	for _, b := range f.Blocks {
		x, err := b.Render(linebreak)
		if err != nil {
			return nil, err
		}
		add(x)
	}

	n := f.Name
	return &file.File{Path: f.Path, Name: n, Mode: filesystem.DefaultMode, Content: util.StringJoin(content, linebreak)}, nil
}
