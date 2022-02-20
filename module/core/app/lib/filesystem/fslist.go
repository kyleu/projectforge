package filesystem

import (
	"os"
	"path/filepath"
	"sort"
	"strings"
)

var defaultIgnore = []string{".DS_Store$", "^.git/", "^.idea/", "^build/", "^client/node_modules", ".html.go$", ".sql.go$"}

func (f *FileSystem) ListFiles(path string, ign []string) []os.DirEntry {
	ignore := buildIgnore(ign)
	infos, err := os.ReadDir(filepath.Join(f.root, path))
	if err != nil {
		f.logger.Warnf("cannot list files in path [%s]: %+v", path, err)
	}
	ret := make([]os.DirEntry, 0, len(infos))
	for _, info := range infos {
		if !checkIgnore(ignore, info.Name()) {
			ret = append(ret, info)
		}
	}
	return ret
}

func (f *FileSystem) ListJSON(path string, trimExtension bool) []string {
	return f.ListExtension(path, "json", trimExtension)
}

func (f *FileSystem) ListExtension(path string, ext string, trimExtension bool) []string {
	matches, err := filepath.Glob(f.getPath(path, "*."+ext))
	if err != nil {
		f.logger.Warnf("cannot list [%s] in path [%s]: %+v", ext, path, err)
	}
	ret := make([]string, 0, len(matches))
	for _, j := range matches {
		idx := strings.LastIndex(j, "/")
		if idx > 0 {
			j = j[idx+1:]
		}
		if trimExtension {
			j = strings.TrimSuffix(j, "."+ext)
		}
		ret = append(ret, j)
	}
	sort.Strings(ret)
	return ret
}

func (f *FileSystem) ListDirectories(path string) []string {
	if !f.Exists(path) {
		return nil
	}
	p := f.getPath(path)
	files, err := os.ReadDir(p)
	if err != nil {
		f.logger.Warnf("cannot list path [%s]: %+v", path, err)
	}
	var ret []string
	for _, f := range files {
		if f.IsDir() {
			ret = append(ret, f.Name())
		}
	}
	sort.Strings(ret)
	return ret
}

func (f *FileSystem) ListFilesRecursive(path string, ign []string) ([]string, error) {
	ignore := buildIgnore(ign)
	p := f.getPath(path)
	var ret []string
	err := filepath.Walk(p, func(fp string, info os.FileInfo, err error) error {
		m := strings.TrimPrefix(fp, p+"/")
		if checkIgnore(ignore, m) {
			return nil
		}
		if !info.IsDir() && strings.Contains(fp, "/") {
			ret = append(ret, m)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	sort.Strings(ret)
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
			if strings.HasPrefix(fp, strings.TrimPrefix(i, keyPrefix)) {
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
