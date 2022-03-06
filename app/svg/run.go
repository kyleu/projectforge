package svg

import (
	"fmt"
	"path/filepath"
	"sort"
	"strings"

	"projectforge.dev/app/lib/filesystem"
)

func Build(fs filesystem.FileLoader) (int, error) {
	return Run(fs, "client/src/svg", "app/util/svg.go")
}

func Run(fs filesystem.FileLoader, src string, tgt string) (int, error) {
	svgs, err := loadSVGs(fs)
	if err != nil {
		return 0, err
	}

	out := template(src, svgs)

	err = writeFile(fs, tgt, out)
	if err != nil {
		return 0, err
	}

	return len(svgs), nil
}

func markup(key string, bytes []byte) string {
	orig := strings.TrimSpace(string(bytes))
	if !strings.Contains(orig, "id=\"svg-") {
		panic(fmt.Sprintf("no id for SVG [%s]", key))
	}
	for strings.Contains(orig, "<!--") {
		startIdx := strings.Index(orig, "<!--")
		endIdx := strings.Index(orig, "-->")
		if endIdx == -1 {
			break
		}
		orig = orig[:startIdx] + orig[endIdx+3:]
	}
	replaced := re.ReplaceAllLiteralString(orig, "")
	return replaced
}

func loadSVGs(fs filesystem.FileLoader) ([]*SVG, error) {
	files := fs.ListExtension(svgPath, "svg", false)
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

	sort.Slice(svgs, func(i int, j int) bool {
		return svgs[i].Key < svgs[j].Key
	})

	return svgs, nil
}

func writeFile(fs filesystem.FileLoader, fn string, out string) error {
	return fs.WriteFile(fn, []byte(out), filesystem.DefaultMode, true)
}
