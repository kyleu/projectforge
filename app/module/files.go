package module

import (
	"os"

	"github.com/pkg/errors"
	"golang.org/x/exp/slices"

	"projectforge.dev/projectforge/app/file"
	"projectforge.dev/projectforge/app/lib/filesystem"
	"projectforge.dev/projectforge/app/util"
)

func (s *Service) GetFilenames(mods Modules, logger util.Logger) ([]string, error) {
	ret := map[string]bool{}
	for _, mod := range mods {
		loader := s.GetFilesystem(mod.Key)
		fs, err := loader.ListFilesRecursive("", nil, logger)
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
	slices.Sort(keys)
	return keys, nil
}

func (s *Service) GetFiles(mods Modules, addHeader bool, logger util.Logger) (file.Files, error) {
	ret := map[string]*file.File{}
	for _, mod := range mods {
		err := s.loadFiles(mod, addHeader, ret, logger)
		if err != nil {
			return nil, errors.Wrapf(err, "unable to load module [%s]", mod.Key)
		}
	}
	keys := make([]string, 0, len(ret))
	for k := range ret {
		keys = append(keys, k)
	}
	slices.Sort(keys)

	files := make(file.Files, 0, len(ret))
	for _, k := range keys {
		files = append(files, ret[k])
	}
	return files, nil
}

func (s *Service) loadFiles(mod *Module, addHeader bool, ret map[string]*file.File, logger util.Logger) error {
	loader := s.GetFilesystem(mod.Key)
	fs, err := loader.ListFilesRecursive("", nil, logger)
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
		fl := file.NewFile(f, mode, b, addHeader, logger)
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
