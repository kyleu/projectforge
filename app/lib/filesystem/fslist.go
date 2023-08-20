// Content managed by Project Forge, see [projectforge.md] for details.
package filesystem

import (
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/samber/lo"

	"projectforge.dev/projectforge/app/util"
)

var defaultIgnore = []string{".DS_Store$", "^.git/", "^.idea/", "^build/", "^client/node_modules", ".html.go$", ".sql.go$"}

func (f *FileSystem) ListFiles(path string, ign []string, logger util.Logger) []fs.DirEntry {
	ignore := buildIgnore(ign)
	infos, err := os.ReadDir(filepath.Join(f.root, path))
	if err != nil {
		logger.Warnf("cannot list files in path [%s]: %+v", path, err)
	}
	return lo.Reject(infos, func(info fs.DirEntry, _ int) bool {
		return checkIgnore(ignore, info.Name())
	})
}

func (f *FileSystem) ListJSON(path string, ign []string, trimExtension bool, logger util.Logger) []string {
	return f.ListExtension(path, "json", ign, trimExtension, logger)
}

func (f *FileSystem) ListExtension(path string, ext string, ign []string, trimExtension bool, logger util.Logger) []string {
	ignore := buildIgnore(ign)
	matches, err := filepath.Glob(f.getPath(path, "*."+ext))
	if err != nil {
		logger.Warnf("cannot list [%s] in path [%s]: %+v", ext, path, err)
	}
	ret := make([]string, 0, len(matches))
	lo.ForEach(matches, func(j string, _ int) {
		if !checkIgnore(ignore, j) {
			idx := strings.LastIndex(j, "/")
			if idx == -1 {
				idx = strings.LastIndex(j, "\\")
			}
			if idx > 0 {
				j = j[idx+1:]
			}
			if trimExtension {
				j = strings.TrimSuffix(j, "."+ext)
			}
			ret = append(ret, j)
		}
	})
	return util.ArraySorted(ret)
}

func (f *FileSystem) ListDirectories(path string, ign []string, logger util.Logger) []string {
	ignore := buildIgnore(ign)
	if !f.Exists(path) {
		return nil
	}
	p := f.getPath(path)
	files, err := os.ReadDir(p)
	if err != nil {
		logger.Warnf("cannot list path [%s]: %+v", path, err)
	}
	ret := lo.FilterMap(files, func(f fs.DirEntry, _ int) (string, bool) {
		if f.IsDir() && !checkIgnore(ignore, f.Name()) {
			return f.Name(), true
		}
		return "", false
	})
	return util.ArraySorted(ret)
}

func (f *FileSystem) ListFilesRecursive(path string, ign []string, _ util.Logger) ([]string, error) {
	ignore := buildIgnore(ign)
	p := f.getPath(path)
	var ret []string
	err := filepath.Walk(p, func(fp string, info fs.FileInfo, err error) error {
		m := strings.TrimPrefix(strings.TrimPrefix(fp, p+"\\"), p+"/")
		if checkIgnore(ignore, m) {
			return nil
		}
		if info != nil && (!info.IsDir()) && (strings.Contains(fp, "/") || strings.Contains(fp, "\\")) {
			ret = append(ret, m)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return util.ArraySorted(ret), nil
}

func (f *FileSystem) Walk(path string, ign []string, fn func(fp string, info fs.FileInfo, err error) error) error {
	ignore := buildIgnore(ign)
	p := f.getPath(path)
	err := filepath.Walk(p, func(fp string, info fs.FileInfo, err error) error {
		m := strings.TrimPrefix(strings.TrimPrefix(fp, p+"\\"), p+"/")
		if checkIgnore(ignore, m) {
			return nil
		}
		return fn(fp, info, err)
	})
	return err
}

func buildIgnore(ign []string) []string {
	ret := append([]string{}, defaultIgnore...)
	ret = append(ret, ign...)
	return ret
}

const (
	keyPrefix = "^"
	keySuffix = "$"
)

func checkIgnore(ignore []string, fp string) bool {
	for _, i := range ignore {
		switch {
		case strings.HasPrefix(i, keyPrefix):
			i = strings.TrimPrefix(i, keyPrefix)
			if fp == strings.TrimSuffix(i, "/") || fp == strings.TrimSuffix(i, "\\") {
				return true
			}
			if strings.HasPrefix(fp, i) {
				return true
			}
		case strings.HasSuffix(i, keySuffix):
			if strings.HasSuffix(fp, strings.TrimSuffix(i, keySuffix)) {
				return true
			}
		case fp == i:
			return true
		}
	}
	return false
}
