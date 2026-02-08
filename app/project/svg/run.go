package svg

import (
	"cmp"
	"slices"
	"strings"

	"github.com/pkg/errors"

	"projectforge.dev/projectforge/app/lib/filesystem"
	"projectforge.dev/projectforge/app/project"
	"projectforge.dev/projectforge/app/util"
)

func Build(fs filesystem.FileLoader, prj *project.Project, logger util.Logger) (int, error) {
	tgt := "app/util/svg.go"
	return Run(fs, tgt, prj, logger)
}

func Run(fs filesystem.FileLoader, tgt string, prj *project.Project, logger util.Logger) (int, error) {
	svgs, err := loadSVGs(prj.Key, fs, logger)
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

func loadSVGs(prj string, fs filesystem.FileLoader, logger util.Logger) (SVGs, error) {
	files, _ := List(fs, logger)
	svgs := make(SVGs, 0, len(files))
	for _, f := range files {
		s, err := Content(prj, fs, f)
		if err != nil {
			return nil, err
		}
		key := strings.TrimSuffix(f, util.ExtSVG)
		mk, err := markup(key, []byte(s))
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
