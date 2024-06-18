package action

import (
	"path"
	"slices"

	"github.com/pkg/errors"

	"projectforge.dev/projectforge/app/lib/filesystem"
	"projectforge.dev/projectforge/app/util"
)

func getEmptyFolders(tgt filesystem.FileLoader, ignore []string, logger util.Logger, pth ...string) ([]string, error) {
	ret := &util.StringSlice{}
	pStr := path.Join(pth...)
	fc := len(tgt.ListFiles(pStr, nil, logger))
	ds := tgt.ListDirectories(pStr, ignore, logger)
	if fc == 0 && len(ds) == 0 {
		ret.Push(pStr)
	}
	for _, d := range ds {
		p := append(slices.Clone(pth), d)
		childRes, err := getEmptyFolders(tgt, ignore, logger, p...)
		if err != nil {
			return nil, errors.Wrapf(err, "unable to get empty folders for [%s/%s]", path.Join(p...), d)
		}
		ret.Push(childRes...)
	}
	return ret.Slice, nil
}
