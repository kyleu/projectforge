package build

import (
	"strings"

	"github.com/pkg/errors"
	"projectforge.dev/projectforge/app/diff"
	"projectforge.dev/projectforge/app/lib/filesystem"
)

func Cleanup(fs filesystem.FileLoader) ([]string, diff.Diffs, error) {
	var logs []string
	var ret diff.Diffs

	files, err := fs.ListFilesRecursive(".", nil)
	if err != nil {
		return nil, nil, errors.Wrap(err, "unable to list project files")
	}

	for _, fn := range files {
		if matches(fn) {
			st, err := fs.Stat(fn)
			if err != nil {
				return nil, nil, errors.Wrapf(err, "can't stat file [%s]", fn)
			}
			if st.Mode() != filesystem.DefaultMode {
				ret = append(ret, &diff.Diff{Path: fn, Status: diff.StatusDifferent, Patch: "fixed mode"})
				err = fs.SetMode(fn, filesystem.DefaultMode)
				if err != nil {
					return nil, nil, errors.Wrapf(err, "can't set mode for file [%s]", fn)
				}
			}
		}
	}

	return logs, ret, nil
}

func matches(fn string) bool {
	return strings.HasSuffix(fn, ".go") || strings.HasSuffix(fn, ".svg") ||
			strings.HasSuffix(fn, ".mod") || strings.HasSuffix(fn, ".json") ||
			strings.HasSuffix(fn, ".html") || strings.HasSuffix(fn, ".sql")
}
