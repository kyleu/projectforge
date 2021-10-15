package filesystem

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

var defaultIgnoreExact = []string{".DS_Store", ".git", ".idea", "build"}
var defaultIgnoreSuffix = []string{".html.go", ".sql.go"}

func (f *FileSystem) ListFiles(path string, ignExact []string, ignSuffix []string) []os.FileInfo {
	ignoreExact, ignoreSuffix := buildIgnore(ignExact, ignSuffix)
	infos, err := ioutil.ReadDir(filepath.Join(f.root, path))
	if err != nil {
		f.logger.Warnf("cannot list files in path [%s]: %+v", path, err)
	}
	ret := make([]os.FileInfo, 0, len(infos))
	for _, info := range infos {
		if !checkIgnore(ignoreExact, ignoreSuffix, info.Name()) {
			ret = append(ret, info)
		}
	}
	return ret
}

func (f *FileSystem) ListJSON(path string, trimExtension bool) []string {
	return f.ListExtension(path, "json", trimExtension)
}

func (f *FileSystem) ListExtension(path string, ext string, trimExtension bool) []string {
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
	sort.Strings(ret)
	return ret
}

func (f *FileSystem) ListFilesRecursive(path string, ignExact []string, ignSuffix []string) ([]string, error) {
	ignoreExact, ignoreSuffix := buildIgnore(ignExact, ignSuffix)
	p := f.getPath(path)
	var ret []string
	err := filepath.Walk(p, func(fp string, info os.FileInfo, err error) error {
		m := strings.TrimPrefix(fp, p+"/")
		if checkIgnore(ignoreExact, ignoreSuffix, m) {
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

func buildIgnore(exact []string, suffix []string) ([]string, []string) {
	e := append([]string{}, defaultIgnoreExact...)
	e = append(e, exact...)
	s := append([]string{}, defaultIgnoreSuffix...)
	s = append(s, suffix...)
	return e, s
}

func checkIgnore(ignoreExact []string, ignoreSuffix []string, fp string) bool {
	fn := fp
	if fni := strings.LastIndex(fp, "/"); fni > -1 {
		fn = fp[fni+1:]
	}
	for _, i := range ignoreExact {
		if fn == i {
			return true
		}
	}
	for _, i := range ignoreSuffix {
		if strings.HasSuffix(fp, i) {
			return true
		}
	}
	return false
}
