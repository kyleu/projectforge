package filesystem

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"{{{ .Package }}}/app/util"
)

var defaultIgnore = []string{".git", ".DS_Store"}

func (f *FileSystem) ListFiles(path string, ignore []string) []os.FileInfo {
	if ignore == nil {
		ignore = defaultIgnore
	}
	infos, err := ioutil.ReadDir(filepath.Join(f.root, path))
	if err != nil {
		f.logger.Warnf("cannot list files in path [%s]: %+v", path, err)
	}
	ret := make([]os.FileInfo, 0, len(infos))
	for _, info := range infos {
		ignored := false
		for _, i := range ignore {
			if i == info.Name() {
				ignored = true
				break
			}
		}
		if !ignored {
			ret = append(ret, info)
		}
	}

	for _, c := range f.children {
		kids := c.ListFiles(path, ignore)
		for _, kid := range kids {
			present := false
			for _, match := range ret {
				if match.Name() == kid.Name() {
					present = true
				}
			}
			if !present {
				ret = append(ret, kid)
			}
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

	for _, c := range f.children {
		kids := c.ListExtension(path, ext, trimExtension)
		for _, kid := range kids {
			if !util.StringArrayContains(ret, kid) {
				ret = append(ret, kid)
			}
		}
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

	for _, c := range f.children {
		kids := c.ListDirectories(path)
		for _, kid := range kids {
			if !util.StringArrayContains(ret, kid) {
				ret = append(ret, kid)
			}
		}
	}

	sort.Strings(ret)
	return ret
}

func (f *FileSystem) ListFilesRecursive(path string, ignore []string) ([]string, error) {
	if ignore == nil {
		ignore = defaultIgnore
	}
	p := f.getPath(path)
	var ret []string
	err := filepath.Walk(p, func(fp string, info os.FileInfo, err error) error {
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

	for _, c := range f.children {
		kids, err := c.ListFilesRecursive(path, ignore)
		if err != nil {
			return nil, err
		}
		for _, kid := range kids {
			if !util.StringArrayContains(ret, kid) {
				ret = append(ret, kid)
			}
		}
	}

	sort.Strings(ret)
	return ret, nil
}
