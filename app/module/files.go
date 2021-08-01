package module

import (
	"os"

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

func (s *Service) LoadAll(keys ...string) (Modules, error) {
	var ret Modules
	for _, key := range keys {
		m, err := s.Load(key)
		if err != nil {
			return nil, err
		}
		ret = append(ret, m)
	}
	return ret, nil
}

func (s *Service) Load(key string) (*Module, error) {
	fs := s.GetFilesystem(key)
	return s.load(key, fs)
}

func (s *Service) load(key string, fs filesystem.FileLoader) (*Module, error) {
	if !fs.Exists(configFilename) {
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

func (s *Service) GetFiles(mods Modules, addHeader bool, tgt filesystem.FileLoader) (file.Files, error) {
	loader := s.GetNestedFilesystem(mods)
	fs, err := loader.ListFilesRecursive("", nil)
	if err != nil {
		return nil, err
	}
	ret := make(file.Files, 0, len(fs))
	for _, f := range fs {
		if f == configFilename {
			continue
		}
		mode, b, err := fileContent(loader, f)
		if err != nil {
			return nil, err
		}
		fl := file.NewFile(f, mode, b, addHeader, s.logger)
		inh, err := file.InheritanceContent(fl)
		if err != nil {
			return nil, err
		}
		if inh != nil {
			err = applyInheritance(fl, inh, addHeader, loader, s.logger)
			if err != nil {
				return nil, err
			}
		}
		err = file.ReplaceSections(fl, tgt)
		if err != nil {
			return nil, err
		}
		ret = append(ret, fl)
	}
	return ret, nil
}
