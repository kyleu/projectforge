package svg

import (
	"strings"

	"github.com/pkg/errors"
	"github.com/samber/lo"

	"projectforge.dev/projectforge/app/lib/filesystem"
	"projectforge.dev/projectforge/app/util"
)

func List(fs filesystem.FileLoader, logger util.Logger) ([]string, error) {
	files := fs.ListExtension("client/src/svg", "svg", nil, true, logger)
	return lo.Map(files, func(key string, _ int) string {
		return strings.TrimSuffix(key, util.ExtSVG)
	}), nil
}

func Content(prj string, fs filesystem.FileLoader, key string) (string, error) {
	c, err := fs.ReadFile(svgPath(prj, key))
	if err != nil {
		return "", errors.Wrapf(err, "unable to read svg file [%s]", key)
	}
	return string(c), nil
}

func Remove(prj string, fs filesystem.FileLoader, key string, logger util.Logger) error {
	return fs.Remove(svgPath(prj, key), logger)
}

func Contents(key string, fs filesystem.FileLoader, logger util.Logger) ([]string, map[string]string, error) {
	files, err := List(fs, logger)
	if err != nil {
		return nil, nil, err
	}
	ret := make(map[string]string, len(files))
	for _, f := range files {
		c, err := Content(key, fs, f)
		if err != nil {
			return nil, nil, err
		}
		ret[f] = c
	}
	return files, ret, nil
}
