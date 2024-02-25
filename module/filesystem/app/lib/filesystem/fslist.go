package filesystem

import (
	"io/fs"
	"path/filepath"
	"strings"

	"github.com/samber/lo"

	"{{{ .Package }}}/app/util"
)

func (f *FileSystem) ListFiles(path string, ign []string, logger util.Logger) FileInfos {
	ignore := buildIgnore(ign)
	p := f.getPath(path)
	dir, err := f.f.Open(p)
	if err != nil {
		if logger != nil {
			logger.Warnf("unable to open path [%s]", p)
		}
		return nil
	}

	infos, err := dir.Readdir(0)
	if err != nil && logger != nil {
		logger.Warnf("cannot list files in path [%s]: %+v", path, err)
	}
	ret := FileInfosFromFS(infos)
	ret = lo.Reject(ret, func(info *FileInfo, _ int) bool {
		return checkIgnore(ignore, info.Name)
	})
	return ret.Sorted()
}

func (f *FileSystem) ListJSON(path string, ign []string, trimExtension bool, logger util.Logger) []string {
	return f.ListExtension(path, "json", ign, trimExtension, logger)
}

func (f *FileSystem) ListExtension(path string, ext string, ign []string, trimExtension bool, logger util.Logger) []string {
	ret := lo.Filter(f.ListFiles(path, ign, logger), func(f *FileInfo, _ int) bool {
		return strings.HasSuffix(f.Name, "."+ext)
	})
	return lo.Map(ret, func(x *FileInfo, _ int) string {
		if trimExtension {
			return strings.TrimSuffix(x.Name, "."+ext)
		}
		return x.Name
	})
}

func (f *FileSystem) ListDirectories(path string, ign []string, logger util.Logger) []string {
	ret := lo.Filter(f.ListFiles(path, ign, logger), func(f *FileInfo, _ int) bool {
		return f.IsDir
	})
	return lo.Map(ret, func(x *FileInfo, _ int) string {
		return x.Name
	})
}

func (f *FileSystem) ListFilesRecursive(path string, ign []string, _ util.Logger) ([]string, error) {
	ignore := buildIgnore(ign)
	p := f.getPath(path)
	var ret []string
	err := filepath.Walk(p, func(fp string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}
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

func (f *FileSystem) Walk(path string, ign []string, fn func(fp string, info *FileInfo, err error) error) error {
	ignore := buildIgnore(ign)
	p := f.getPath(path)
	err := filepath.Walk(p, func(fp string, info fs.FileInfo, err error) error {
		m := strings.TrimPrefix(strings.TrimPrefix(fp, p+"\\"), p+"/")
		if checkIgnore(ignore, m) {
			return nil
		}
		return fn(fp, FileInfoFromFS(info), err)
	})
	return err
}
