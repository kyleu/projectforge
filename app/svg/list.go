package svg

import (
	"path/filepath"
	"strings"

	"projectforge.dev/app/lib/filesystem"
	"github.com/pkg/errors"
)

func List(fs filesystem.FileLoader) ([]string, error) {
	files := fs.ListExtension(svgPath, "svg", false)
	ret := make([]string, 0, len(files))
	for _, key := range files {
		ret = append(ret, strings.TrimSuffix(key, ".svg"))
	}
	return ret, nil
}

func Content(fs filesystem.FileLoader, key string) (string, error) {
	c, err := fs.ReadFile(filepath.Join(svgPath, key+".svg"))
	if err != nil {
		return "", errors.Wrapf(err, "unable to read svg file [%s]", key)
	}
	return string(c), nil
}

func Remove(fs filesystem.FileLoader, key string) error {
	return fs.Remove(filepath.Join(svgPath, key+".svg"))
}

func Contents(fs filesystem.FileLoader) ([]string, map[string]string, error) {
	files, err := List(fs)
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
