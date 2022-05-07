package build

import (
	"fmt"
	"strings"

	"github.com/pkg/errors"
	"projectforge.dev/projectforge/app/diff"
	"projectforge.dev/projectforge/app/lib/filesystem"
)

func Imports(self string, fix bool, targetPath string, fs filesystem.FileLoader) ([]string, diff.Diffs, error) {
	var logs []string
	var ret diff.Diffs

	files, err := fs.ListFilesRecursive(".", nil)
	if err != nil {
		return nil, nil, errors.Wrap(err, "unable to list project files")
	}

	for _, fn := range files {
		if !strings.HasSuffix(fn, ".go") {
			continue
		}
		content, err := fs.ReadFile(fn)
		if err != nil {
			return nil, nil, errors.Wrapf(err, "unable to read file [%s]", fn)
		}
		stat, err := fs.Stat(fn)
		if err != nil {
			return nil, nil, err
		}
		fixed, diffs, err := processFile(fn, strings.Split(string(content), "\n"), self)
		ret = append(ret, diffs...)
		if err != nil {
			return nil, nil, errors.Wrapf(err, "unable to process imports for [%s]", fn)
		}
		if fix && len(diffs) > 0 {
			if targetPath == "" || fn == targetPath {
				newContent := strings.Join(fixed, "\n")
				err = fs.WriteFile(fn, []byte(newContent), stat.Mode(), true)
				if err != nil {
					return nil, nil, errors.Wrapf(err, "unable to write file [%s]", fn)
				}
				logs = append(logs, fmt.Sprintf("fixed imports for [%s]", fn))
			}
		}
	}

	return logs, ret, nil
}

func processFile(fn string, lines []string, self string) ([]string, diff.Diffs, error) {
	if strings.HasPrefix(fn, "module/") {
		return nil, nil, nil
	}
	var ret diff.Diffs

	var started bool
	var start int
	var end int
	var imports, orig []string
	for idx, line := range lines {
		if started {
			if line == ")" {
				if end > 0 {
					return nil, nil, errors.New("multiple import section endings")
				}
				end = idx
				break
			}
			i := strings.TrimSpace(line)
			if lidx := strings.Index(line, " "); lidx > -1 {
				i = strings.TrimSpace(line[lidx:])
			}
			i = strings.TrimPrefix(i, "_ ")
			i = strings.TrimSuffix(i, "\"")
			i = strings.TrimPrefix(i, "\"")
			imports = append(imports, i+":"+importType(i, self))
			orig = append(orig, line)
		} else {
			if strings.HasPrefix(line, "import") && strings.Contains(line, "(") {
				started = true
				if start > 0 {
					return nil, nil, errors.New("multiple import section starts")
				}
				start = idx
			}
			continue
		}
	}

	ordered, chkErr := check(imports, orig)
	if chkErr == nil {
		return lines, ret, nil
	}

	ret = append(ret, &diff.Diff{
		Path:    fn,
		Status:  diff.StatusDifferent,
		Patch:   fmt.Sprintf("%s: %s", fn, chkErr.Error()),
		Changes: nil,
	})

	fixed := make([]string, 0, len(lines))
	fixed = append(fixed, lines[:start+1]...)
	fixed = append(fixed, ordered...)
	fixed = append(fixed, lines[end:]...)

	return fixed, ret, nil
}
