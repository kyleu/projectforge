package svg

import (
	"cmp"
	"path/filepath"
	"slices"
	"strings"

	"github.com/pkg/errors"

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

	if len(svgs) == 0 {
		return 0, errors.New("no SVGs available")
	}

	out := template(svgs, util.StringDetectLinebreak(svgs[0].Markup))

	err = writeFile(fs, tgt, out)
	if err != nil {
		return 0, err
	}

	return len(svgs), nil
}

func markup(key string, bytes []byte) (string, error) {
	orig, _ := cleanMarkup(strings.TrimSpace(string(bytes)), "")
	if !strings.Contains(orig, "id=\"svg-") {
		return "", errors.Errorf("no id for SVG [%s]", key)
	}
	replaced := re.ReplaceAllLiteralString(orig, "")
	return replaced, nil
}

func loadSVGs(fs filesystem.FileLoader, logger util.Logger) ([]*SVG, error) {
	files := fs.ListExtension(svgPath, "svg", nil, false, logger)
	svgs := make([]*SVG, 0, len(files))
	for _, f := range files {
		b, err := fs.ReadFile(filepath.Join(svgPath, f))
		if err != nil {
			return nil, err
		}
		key := strings.TrimSuffix(f, util.ExtSVG)
		mk, err := markup(key, b)
		if err != nil {
			return nil, err
		}
		svgs = append(svgs, &SVG{Key: key, Markup: mk})
	}

	slices.SortFunc(svgs, func(l *SVG, r *SVG) int {
		return cmp.Compare(l.Key, r.Key)
	})

	return svgs, nil
}

func writeFile(fs filesystem.FileLoader, fn string, out string) error {
	return fs.WriteFile(fn, []byte(out), filesystem.DefaultMode, true)
}
