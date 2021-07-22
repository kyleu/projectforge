package svg

import (
	"github.com/kyleu/projectforge/app/filesystem"
)

func List(fs filesystem.FileLoader) ([]string, error) {
	files := fs.ListExtension(svgPath, "svg", false)
	return files, nil
}
