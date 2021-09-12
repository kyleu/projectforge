package module

import (
	"os"
	"sort"

	"github.com/kyleu/projectforge/app/file"
	"github.com/kyleu/projectforge/app/filesystem"
	"github.com/kyleu/projectforge/app/util"
	"github.com/pkg/errors"
)

const configFilename = ".module.json"

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

func (s *Service) LoadNative(keys ...string) (Modules, error) {
	var ret Modules
	for _, key := range keys {
		m, err := s.Load(key, "")
		if err != nil {
			return nil, err
		}
		ret = append(ret, m)
	}
	return ret, nil
}

func (s *Service) Load(key string, url string) (*Module, error) {
	fs := s.GetFilesystem(key)
	return s.load(key, fs, url)
}

func (s *Service) load(key string, fs filesystem.FileLoader, url string) (*Module, error) {
	if !fs.Exists(configFilename) {
		// TODO load from URL
		msg := "file [%s] does not exist in path for module [%s] using root [%s]"
		return nil, errors.Errorf(msg, configFilename, key, fs.Root())
	}

	b, err := fs.ReadFile(configFilename)
	if err != nil {
		return nil, errors.Wrapf(err, "unable to read [%s]", configFilename)
	}

	ret := &Module{}
	err = util.FromJSON(b, ret)
	if err != nil {
		return nil, errors.Wrapf(err, "unable to parse [%s] as module with key [%s]", configFilename, key)
	}
	ret.Key = key
	ret.Files = fs

	s.cache[key] = ret
	return ret, nil
}

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

func (s *Service) GetFiles(mods Modules, addHeader bool, tgt filesystem.FileLoader) (file.Files, error) {
	ret := map[string]*file.File{}
	for _, mod := range mods {
		loader := s.GetFilesystem(mod.Key)
		fs, err := loader.ListFilesRecursive("", nil)
		if err != nil {
			return nil, err
		}
		for _, f := range fs {
			if f == configFilename {
				continue
			}
			mode, b, err := fileContent(loader, f)
			if err != nil {
				return nil, err
			}
			fl := file.NewFile(f, mode, b, addHeader, s.logger)

			curr, exists := ret[fl.FullPath()]
			if exists {
				inh, err := file.InheritanceContent(fl)
				if err != nil {
					return nil, err
				}
				if inh != nil {
					err = applyInheritance(fl, inh, curr.Content)
					if err != nil {
						return nil, err
					}
				}
			}
			ret[fl.FullPath()] = fl
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
