package module

import (
	"context"
	"path/filepath"

	"github.com/pkg/errors"
	"projectforge.dev/app/lib/filesystem"
	"projectforge.dev/app/util"
)

const configFilename = ".module.json"

func (s *Service) LoadNative(ctx context.Context, keys ...string) (Modules, error) {
	var ret Modules
	for _, key := range keys {
		m, err := s.Load(ctx, key, "")
		if err != nil {
			return nil, err
		}
		ret = append(ret, m)
	}
	return ret, nil
}

func (s *Service) Load(ctx context.Context, key string, url string) (*Module, error) {
	return s.load(ctx, key, "", url)
}

func (s *Service) load(ctx context.Context, key string, path string, url string) (*Module, error) {
	var fs filesystem.FileLoader
	switch {
	case s.local.IsDir(key):
		fs = filesystem.NewFileSystem(filepath.Join(s.local.Root(), key), s.logger)
	case s.config.IsDir(key):
		fs = filesystem.NewFileSystem(filepath.Join(s.config.Root(), key), s.logger)
	case path != "":
		fs = filesystem.NewFileSystem(path, s.logger)
	default:
		var err error
		if url == "" {
			url, err = s.AssetURL(ctx, key)
			if err != nil {
				return nil, errors.Wrapf(err, "error downloading module [%s]", key)
			}
		}
		err = s.Download(key, url)
		if err != nil {
			return nil, errors.Wrapf(err, "error downloading module [%s]", key)
		}
		fs = filesystem.NewFileSystem(filepath.Join(s.config.Root(), key), s.logger)
	}
	if !fs.Exists(configFilename) {
		const msg = "file [%s] does not exist in path for module [%s] using root [%s]"
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

	if fs.Exists(ret.DocPath()) {
		b, err = fs.ReadFile(ret.DocPath())
		if err != nil {
			return nil, errors.Wrapf(err, "unable to read [%s]", configFilename)
		}
		ret.UsageMD = string(b)
	}
	s.cache[key] = ret
	return ret, nil
}
