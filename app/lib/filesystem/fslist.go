// Content managed by Project Forge, see [projectforge.md] for details.
package filesystem

import (
	"os"
	"path/filepath"
	"strings"

	"golang.org/x/exp/slices"

	"projectforge.dev/projectforge/app/util"
)

var defaultIgnore = []string{".DS_Store$", "^.git/", "^.idea/", "^build/", "^client/node_modules", ".html.go$", ".sql.go$"}

func (f *FileSystem) ListFiles(path string, ign []string, logger util.Logger) []os.DirEntry {
	ignore := buildIgnore(ign)
	infos, err := os.ReadDir(filepath.Join(f.root, path))
	if err != nil {
		logger.Warnf("cannot list files in path [%s]: %+v", path, err)
	}
	ret := make([]os.DirEntry, 0, len(infos))
	for _, info := range infos {
		if !checkIgnore(ignore, info.Name()) {
			ret = append(ret, info)
		}
	}
	return ret
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
	for _, j := range matches {
		if !checkIgnore(ignore, j) {
			idx := strings.LastIndex(j, "/")
			if idx > 0 {
				j = j[idx+1:]
			}
			if trimExtension {
				j = strings.TrimSuffix(j, "."+ext)
			}
			ret = append(ret, j)
		}
	}
	slices.Sort(ret)
	return ret
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
	var ret []string
	for _, f := range files {
		if f.IsDir() {
			if !checkIgnore(ignore, f.Name()) {
				ret = append(ret, f.Name())
			}
		}
	}
	slices.Sort(ret)
	return ret
}

func (f *FileSystem) ListFilesRecursive(path string, ign []string, logger util.Logger) ([]string, error) {
	ignore := buildIgnore(ign)
	p := f.getPath(path)
	var ret []string
	err := filepath.Walk(p, func(fp string, info os.FileInfo, err error) error {
		m := strings.TrimPrefix(fp, p+"/")
		if checkIgnore(ignore, m) {
			return nil
		}
		if info != nil && (!info.IsDir()) && strings.Contains(fp, "/") {
			ret = append(ret, m)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	slices.Sort(ret)
	return ret, nil
}

func (f *FileSystem) Walk(path string, ign []string, fn func(fp string, info os.FileInfo, err error) error) error {
	ignore := buildIgnore(ign)
	p := f.getPath(path)
	err := filepath.Walk(p, func(fp string, info os.FileInfo, err error) error {
		m := strings.TrimPrefix(fp, p+"/")
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
			if strings.HasSuffix(i, "/") && fp == strings.TrimSuffix(i, "/") {
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
