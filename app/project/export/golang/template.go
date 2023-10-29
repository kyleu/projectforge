package golang

import (
	"fmt"
	"strings"

	"github.com/samber/lo"

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
	lo.ForEach(i, func(imp *Import, _ int) {
		hit := lo.ContainsBy(f.Imports, func(x *Import) bool {
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

func (f *Template) Render(addHeader bool, linebreak string) (*file.File, error) {
	var content []string
	add := func(s string, args ...any) {
		content = append(content, fmt.Sprintf(s+linebreak, args...))
	}

	if addHeader {
		switch {
		case strings.HasSuffix(f.Name, ".sql"):
			content = append(content, fmt.Sprintf("-- %s", file.HeaderContent))
		case strings.HasSuffix(f.Name, ".graphql"):
			content = append(content, fmt.Sprintf("# %s", file.HeaderContent))
		default:
			content = append(content, fmt.Sprintf("<!-- %s -->", file.HeaderContent))
		}
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
	return &file.File{Path: f.Path, Name: n, Mode: filesystem.DefaultMode, Content: strings.Join(content, linebreak)}, nil
}
