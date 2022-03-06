package module

import (
	"os"
	"sort"

	"github.com/pkg/errors"
	"projectforge.dev/app/file"
	"projectforge.dev/app/lib/filesystem"
)

func (s *Service) GetFilenames(mods Modules) ([]string, error) {
	ret := map[string]bool{}
	for _, mod := range mods {
		loader := s.GetFilesystem(mod.Key)
		fs, err := loader.ListFilesRecursive("", nil)
		if err != nil {
			return nil, err
		}
		for _, f := range fs {
			ret[f] = true
		}
	}
	keys := make([]string, 0, len(ret))
	for k := range ret {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys, nil
}

func (s *Service) GetFiles(mods Modules, addHeader bool) (file.Files, error) {
	ret := map[string]*file.File{}
	for _, mod := range mods {
		err := s.loadFiles(mod, addHeader, ret)
		if err != nil {
			return nil, errors.Wrapf(err, "unable to load module [%s]", mod.Key)
		}
	}
	keys := make([]string, 0, len(ret))
	for k := range ret {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	files := make(file.Files, 0, len(ret))
	for _, k := range keys {
		files = append(files, ret[k])
	}
	return files, nil
}

func (s *Service) loadFiles(mod *Module, addHeader bool, ret map[string]*file.File) error {
	loader := s.GetFilesystem(mod.Key)
	fs, err := loader.ListFilesRecursive("", nil)
	if err != nil {
		return err
	}
	for _, f := range fs {
		if f == configFilename {
			continue
		}
		mode, b, err := fileContent(loader, f)
		if err != nil {
			return err
		}
		fl := file.NewFile(f, mode, b, addHeader, s.logger)
		ret[fl.FullPath()] = fl
	}
	return nil
}

func fileContent(files filesystem.FileLoader, path string) (os.FileMode, []byte, error) {
	stat, err := files.Stat(path)
	if err != nil {
		return 0, nil, errors.Wrapf(err, "file [%s] not found", path)
	}
	b, err := files.ReadFile(path)
	if err != nil {
		return 0, nil, err
	}
	return stat.Mode(), b, nil
}
