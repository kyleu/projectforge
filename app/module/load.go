package module

import (
	"path/filepath"

	"github.com/kyleu/projectforge/app/filesystem"
	"github.com/kyleu/projectforge/app/util"
	"github.com/pkg/errors"
)

const configFilename = ".module.json"
const summaryFilename = ".module.md"

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
	return s.load(key, "", url)
}

func (s *Service) load(key string, path string, url string) (*Module, error) {
	var fs filesystem.FileLoader
	switch {
	case s.local.IsDir(key):
		fs = filesystem.NewFileSystem(filepath.Join(s.local.Root(), key), s.logger)
	case s.config.IsDir(key):
		fs = filesystem.NewFileSystem(filepath.Join(s.config.Root(), key), s.logger)
	case path != "":
		fs = filesystem.NewFileSystem(path, s.logger)
	default:
		err := s.Download(key, url)
		if err != nil {
			return nil, errors.Wrapf(err, "error downloading module [%s]", key)
		}
		fs = filesystem.NewFileSystem(filepath.Join(s.config.Root(), key), s.logger)
	}
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
	ret.URL = url

	if fs.Exists(summaryFilename) {
		b, err = fs.ReadFile(summaryFilename)
		if err != nil {
			return nil, errors.Wrapf(err, "unable to read [%s]", configFilename)
		}
		ret.UsageMD = string(b)
	}
	s.cache[key] = ret
	return ret, nil
}
