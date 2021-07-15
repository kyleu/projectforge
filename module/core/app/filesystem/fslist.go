package filesystem

import (
	"io/fs"
	"io/ioutil"
	"path/filepath"
	"strings"
)

var defaultIgnore = []string{".git", ".DS_Store"}

func (f *FileSystem) ListFiles(path string) []string {
	infos, err := ioutil.ReadDir(path)
	if err != nil {
		f.logger.Warnf("cannot list files in path [%s]: %+v", path, err)
	}
	matches := make([]string, 0, len(infos))
	for _, info := range infos {
		matches = append(matches, info.Name())
	}
	return matches
}

func (f *FileSystem) ListJSON(path string) []string {
	return f.ListExtension(path, "json")
}

func (f *FileSystem) ListExtension(path string, ext string) []string {
	glob := "*." + ext
	matches, err := filepath.Glob(f.getPath(path, glob))
	if err != nil {
		f.logger.Warnf("cannot list [%s] in path [%s]: %+v", ext, path, err)
	}
	ret := make([]string, 0, len(matches))
	for _, j := range matches {
		idx := strings.LastIndex(j, "/")
		if idx > 0 {
			j = j[idx+1:]
		}
		ret = append(ret, strings.TrimSuffix(j, "."+ext))
	}
	return ret
}

func (f *FileSystem) ListDirectories(path string) []string {
	if !f.Exists(path) {
		return nil
	}
	p := f.getPath(path)
	files, err := ioutil.ReadDir(p)
	if err != nil {
		f.logger.Warnf("cannot list path [%s]: %+v", path, err)
	}
	var ret []string
	for _, f := range files {
		if f.IsDir() {
			ret = append(ret, f.Name())
		}
	}
	return ret
}

func (f *FileSystem) ListFilesRecursive(path string, ignore []string) ([]string, error) {
	if ignore == nil {
		ignore = defaultIgnore
	}
	p := f.getPath(path)
	var ret []string
	err := filepath.Walk(p, func(fp string, info fs.FileInfo, err error) error {
		for _, i := range ignore {
			if strings.HasPrefix(fp, i) || strings.HasSuffix(fp, i) {
				return nil
			}
		}
		if !info.IsDir() {
			ret = append(ret, strings.TrimPrefix(strings.TrimPrefix(fp, p), "/"))
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return ret, nil
}
