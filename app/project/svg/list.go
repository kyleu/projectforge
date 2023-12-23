package svg

import (
	"path/filepath"
	"strings"

	"github.com/pkg/errors"
	"github.com/samber/lo"

	"projectforge.dev/projectforge/app/lib/filesystem"
	"projectforge.dev/projectforge/app/util"
)

func List(fs filesystem.FileLoader, logger util.Logger) ([]string, error) {
	files := fs.ListExtension(svgPath, "svg", nil, false, logger)
	return lo.Map(files, func(key string, _ int) string {
		return strings.TrimSuffix(key, util.ExtSVG)
	}), nil
}

func Content(fs filesystem.FileLoader, key string) (string, error) {
	c, err := fs.ReadFile(filepath.Join(svgPath, key+util.ExtSVG))
	if err != nil {
		return "", errors.Wrapf(err, "unable to read svg file [%s]", key)
	}
	return string(c), nil
}

func Remove(fs filesystem.FileLoader, key string, logger util.Logger) error {
	return fs.Remove(filepath.Join(svgPath, key+util.ExtSVG), logger)
}

func Contents(fs filesystem.FileLoader, logger util.Logger) ([]string, map[string]string, error) {
	files, err := List(fs, logger)
	if err != nil {
		return nil, nil, err
	}
	ret := make(map[string]string, len(files))
	for _, key := range files {
		c, err := Content(fs, key)
		if err != nil {
			return nil, nil, err
		}
		ret[key] = c
	}
	return files, ret, nil
}
