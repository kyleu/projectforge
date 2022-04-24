package build

import (
	"fmt"
	"strings"

	"github.com/pkg/errors"
	"projectforge.dev/projectforge/app/diff"
	"projectforge.dev/projectforge/app/lib/filesystem"
	"projectforge.dev/projectforge/app/util"
)

func Imports(self string, fix bool, fs filesystem.FileLoader) (diff.Diffs, error) {
	var ret diff.Diffs

	files, err := fs.ListFilesRecursive(".", nil)
	if err != nil {
		return nil, errors.Wrap(err, "unable to list project files")
	}

	for _, fn := range files {
		if strings.HasSuffix(fn, ".go") {
			content, err := fs.ReadFile(fn)
			if err != nil {
				return nil, errors.Wrapf(err, "unable to read file [%s]", fn)
			}

			diffs, err := processFile(fn, strings.Split(string(content), "\n"), self)
			ret = append(ret, diffs...)
			if err != nil {
				if fix {
					return nil, errors.Wrapf(err, "unable to process imports for [%s]", fn)
				} else {
					return nil, errors.Wrapf(err, "unable to process imports for [%s]", fn)
				}
			}
		}
	}

	return ret, nil
}

func processFile(fn string, lines []string, self string) (diff.Diffs, error) {
	if strings.HasPrefix(fn, "module/") {
		return nil, nil
	}
	var ret diff.Diffs

	var started bool
	var start int
	var end int
	var imports []string
	for idx, line := range lines {
		if started {
			if line == ")" {
				started = false
				if end > 0 {
					return nil, errors.New("multiple import section endings")
				}
				end = idx
				break
			}
			i := strings.TrimPrefix(line, "_ ")
			i = strings.TrimSpace(i)
			i = strings.TrimSuffix(i, "\"")
			i = strings.TrimPrefix(i, "\"")
			imports = append(imports, i+":"+importType(i, self))
		} else {
			if strings.HasPrefix(line, "import") && strings.Contains(line, "(") {
				started = true
				if start > 0 {
					return nil, errors.New("multiple import section starts")
				}
				start = idx
			}
			continue
		}
	}
	if chkErr := check(imports); chkErr != nil {
		ret = append(ret, &diff.Diff{
			Path:    fn,
			Status:  diff.StatusDifferent,
			Patch:   fmt.Sprintf("%s: %s", fn, chkErr.Error()),
			Changes: nil,
		})
	}
	return ret, nil
}

func check(imports []string) error {
	var state int
	var lastLine string
	var observed []string
	observe := func(key string, i string) error {
		for _, ob := range observed {
			if ob > i {
				return errors.Errorf("%s sorting", key)
			}
		}
		observed = append(observed, i)
		return nil
	}
	clear := func() {
		observed = []string{}
	}

	for _, imp := range imports {
		i, l := util.StringSplitLast(imp, ':', true)
		switch l {
		case "sep":
			if state != 0 && lastLine != "" {
				state++
				clear()
			}
		case "1st":
			if state > 1 {
				return errors.New("1st party")
			}
			if state != 1 {
				state = 1
				clear()
			}
			if err := observe(i, "1st party"); err != nil {
				return err
			}
		case "3rd":
			if state > 2 {
				return errors.New("3rd party")
			}
			if state != 2 {
				state = 2
				clear()
			}
			if err := observe(i, "3rd party"); err != nil {
				return err
			}
		case "self":
			if state > 3 {
				return errors.New("self")
			}
			if state != 3 {
				state = 3
				clear()
			}
			if err := observe(i, "self"); err != nil {
				return err
			}
		default:
			return errors.New("invalid type")
		}
		lastLine = l
	}
	return nil
}

func importType(i string, self string) string {
	if i == "" {
		return "sep"
	}
	if strings.HasPrefix(i, self) || strings.HasPrefix(i, "{{{ .Package }}}") {
		return "self"
	}
	if strings.Contains(i, ".") {
		return "3rd"
	}
	return "1st"
}
