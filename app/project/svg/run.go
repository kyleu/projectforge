package svg

import (
	"fmt"
	"path/filepath"
	"strings"

	"golang.org/x/exp/slices"

	"projectforge.dev/projectforge/app/lib/filesystem"
	"projectforge.dev/projectforge/app/util"
)

func Build(fs filesystem.FileLoader, logger util.Logger) (int, error) {
	return Run(fs, "app/util/svg.go", logger)
}

func Run(fs filesystem.FileLoader, tgt string, logger util.Logger) (int, error) {
	svgs, err := loadSVGs(fs, logger)
	if err != nil {
		return 0, err
	}

	out := template(svgs)

	err = writeFile(fs, tgt, out)
	if err != nil {
		return 0, err
	}

	return len(svgs), nil
}

func markup(key string, bytes []byte) string {
	orig, _ := cleanMarkup(strings.TrimSpace(string(bytes)), "")
	if !strings.Contains(orig, "id=\"svg-") {
		panic(fmt.Sprintf("no id for SVG [%s]", key))
	}
	replaced := re.ReplaceAllLiteralString(orig, "")
	return replaced
}

func loadSVGs(fs filesystem.FileLoader, logger util.Logger) ([]*SVG, error) {
	files := fs.ListExtension(svgPath, "svg", nil, false, logger)
	svgs := make([]*SVG, 0, len(files))
	for _, f := range files {
		b, err := fs.ReadFile(filepath.Join(svgPath, f))
		if err != nil {
			panic(err)
		}
		key := strings.TrimSuffix(f, ".svg")
		svgs = append(svgs, &SVG{
			Key:    key,
			Markup: markup(key, b),
		})
	}

	slices.SortFunc(svgs, func(l *SVG, r *SVG) bool {
		return l.Key < r.Key
	})

	return svgs, nil
}

func writeFile(fs filesystem.FileLoader, fn string, out string) error {
	return fs.WriteFile(fn, []byte(out), filesystem.DefaultMode, true)
}
