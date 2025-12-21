package build

import (
	"fmt"
	"strings"

	"github.com/pkg/errors"
	"github.com/samber/lo"

	"projectforge.dev/projectforge/app/file/diff"
	"projectforge.dev/projectforge/app/lib/filesystem"
	"projectforge.dev/projectforge/app/util"
)

const OSWindows = "windows"

func Imports(self string, fix bool, targetPath string, fs filesystem.FileLoader, logger util.Logger) ([]string, diff.Diffs, error) {
	var logs []string
	var ret diff.Diffs

	files, err := fs.ListFilesRecursive(".", nil, logger)
	if err != nil {
		return nil, nil, errors.Wrap(err, "unable to list project files")
	}

	for _, fn := range files {
		if strings.HasSuffix(fn, "generated.go") || strings.HasSuffix(fn, ".pb.go") || strings.HasSuffix(fn, "schema.resolvers.go") {
			continue
		}
		rlogs, rdiffs, err := importsFor(self, fix, fs, fn, targetPath)
		if err != nil {
			return nil, nil, errors.Wrap(err, "")
		}
		logs = append(logs, rlogs...)
		ret = append(ret, rdiffs...)
	}

	return logs, ret, nil
}

func importsFor(self string, fix bool, fs filesystem.FileLoader, fn string, targetPath string) ([]string, diff.Diffs, error) {
	if !strings.HasSuffix(fn, util.ExtGo) && !strings.HasSuffix(fn, util.ExtHTML) {
		return nil, nil, nil
	}
	var ret diff.Diffs
	var logs []string
	content, err := fs.ReadFile(fn)
	if err != nil {
		return nil, nil, errors.Wrapf(err, "unable to read file [%s]", fn)
	}
	str := string(content)
	stat, err := fs.Stat(fn)
	if err != nil {
		return nil, nil, err
	}
	_, fixed, diffs, err := processFileImports(fn, util.StringSplitLines(str), self)
	ret = append(ret, diffs...)
	if err != nil {
		return nil, nil, errors.Wrapf(err, "unable to process imports for [%s]", fn)
	}
	if fix && len(diffs) > 0 {
		if targetPath == "" || fn == targetPath {
			newContent := util.StringJoin(fixed, util.StringDetectLinebreak(str))
			err = fs.WriteFile(fn, []byte(newContent), stat.Mode, true)
			if err != nil {
				return nil, nil, errors.Wrapf(err, "unable to write file [%s]", fn)
			}
			logs = append(logs, fmt.Sprintf("fixed imports for [%s]", fn))
		}
	}
	return logs, ret, nil
}

func processFileImports(fn string, lines []string, self string) ([]string, []string, diff.Diffs, error) {
	if strings.HasPrefix(fn, "module/") {
		return nil, nil, nil, nil
	}
	var ret diff.Diffs

	var started bool
	var start int
	var end int
	var imports, orig []string
	for idx, line := range lines {
		if started {
			if line == ")" || line == ") %}" {
				if end > 0 {
					return nil, nil, nil, errors.New("multiple import section endings")
				}
				end = idx
				break
			}
			i := strings.TrimSpace(line)
			if lidx := strings.Index(line, " "); lidx > -1 {
				i = strings.TrimSpace(line[lidx:])
			}
			i = strings.TrimPrefix(i, "_ ")
			i = strings.TrimSuffix(i, `"`)
			i = strings.TrimPrefix(i, `"`)
			imports = append(imports, i+":"+importType(i, self))
			orig = append(orig, line)
		} else {
			if (strings.HasPrefix(line, "import") || strings.HasPrefix(line, "{% import")) && strings.Contains(line, "(") {
				started = true
				if start > 0 {
					return nil, nil, nil, errors.New("multiple import section starts")
				}
				start = idx
			}
			continue
		}
	}

	_, ordered, chkErr := check(imports, orig, lines)
	if chkErr == nil {
		return imports, lines, ret, nil
	}

	ret = append(ret, &diff.Diff{Path: fn, Status: diff.StatusDifferent, Patch: fmt.Sprintf("%s: %s", fn, chkErr.Error())})

	funky := lo.ContainsBy(ordered, func(item string) bool {
		return strings.Contains(item, "_") || strings.Contains(item, "//") || strings.Contains(item, "/*")
	})

	fixed := make([]string, 0, len(lines))
	if len(ordered) == 1 && !funky {
		fixed = append(fixed, lines[:start]...)
		if util.RS(lines[start]).Contains("{%") {
			fixed = append(fixed, util.RS(ordered[0]).TrimSpace().WithPrefix("{% import ").WithSuffix(" %}").String())
		} else {
			fixed = append(fixed, util.RS(ordered[0]).TrimSpace().WithPrefix("import ").String())
		}
		fixed = append(fixed, lines[end+1:]...)
	} else {
		fixed = append(fixed, lines[:start]...)
		fixed = append(fixed, "import (")
		fixed = append(fixed, ordered...)
		fixed = append(fixed, ")")
		fixed = append(fixed, lines[end+1:]...)
	}
	return imports, fixed, ret, nil
}
