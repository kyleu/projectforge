package build

import (
	"strings"

	"github.com/pkg/errors"

	"projectforge.dev/projectforge/app/file/diff"
	"projectforge.dev/projectforge/app/lib/filesystem"
	"projectforge.dev/projectforge/app/util"
)

func Cleanup(fs filesystem.FileLoader, logger util.Logger) ([]string, diff.Diffs, error) {
	var logs []string
	var ret diff.Diffs

	files, err := fs.ListFilesRecursive(".", nil, logger)
	if err != nil {
		return nil, nil, errors.Wrap(err, "unable to list project files")
	}

	for _, fn := range files {
		if !matches(fn) {
			continue
		}
		st, err := fs.Stat(fn)
		if err != nil {
			return nil, nil, errors.Wrapf(err, "can't stat file [%s]", fn)
		}
		if st.Mode != filesystem.DefaultMode {
			ret = append(ret, &diff.Diff{Path: fn, Status: diff.StatusDifferent, Patch: "fixed mode"})
			err = fs.SetMode(fn, filesystem.DefaultMode)
			if err != nil {
				return nil, nil, errors.Wrapf(err, "can't set mode for file [%s]", fn)
			}
		}
	}

	return logs, ret, nil
}

func matches(fn string) bool {
	return strings.HasSuffix(fn, util.ExtGo) ||
		strings.HasSuffix(fn, util.ExtJS) || strings.HasSuffix(fn, util.ExtTS) ||
		strings.HasSuffix(fn, util.ExtCSS) || strings.HasSuffix(fn, util.ExtSVG) ||
		strings.HasSuffix(fn, util.ExtMod) || strings.HasSuffix(fn, util.ExtJSON) ||
		strings.HasSuffix(fn, util.ExtHTML) || strings.HasSuffix(fn, util.ExtSQL)
}
