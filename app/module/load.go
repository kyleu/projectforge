package module

import (
	"context"
	"path"
	"path/filepath"

	"github.com/pkg/errors"
	"projectforge.dev/projectforge/app/lib/filesystem"
	"projectforge.dev/projectforge/app/util"
)

const configFilename = ".module.json"

func (s *Service) LoadNative(ctx context.Context, keys ...string) (Modules, error) {
	var ret Modules
	for _, key := range keys {
		m, err := s.Load(ctx, key, "")
		if err != nil {
			return nil, err
		}
		ret = append(ret, m...)
	}
	return ret, nil
}

func (s *Service) Load(ctx context.Context, key string, url string) (Modules, error) {
	return s.load(ctx, key, "", url)
}

func (s *Service) load(ctx context.Context, key string, pth string, url string) (Modules, error) {
	var fs filesystem.FileLoader
	switch {
	case s.local.IsDir(key):
		fs = filesystem.NewFileSystem(filepath.Join(s.local.Root(), key), s.logger)
	case s.config.IsDir(key):
		fs = filesystem.NewFileSystem(filepath.Join(s.config.Root(), key), s.logger)
	case pth != "":
		fs = filesystem.NewFileSystem(pth, s.logger)
		if key == "*" {
			return s.loadDirectory(ctx, pth, url, fs)
		}
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
	s.cacheMu.Lock()
	s.cache[key] = ret
	s.cacheMu.Unlock()
	return Modules{ret}, nil
}

func (s *Service) loadDirectory(ctx context.Context, pth string, u string, fs filesystem.FileLoader) (Modules, error) {
	dirs := fs.ListDirectories(pth, nil)
	if len(dirs) == 0 {
		return nil, errors.Errorf("directory at path [%s] does not contain module directories", pth)
	}
	ret := make(Modules, 0, len(dirs))
	for _, dir := range dirs {
		if !fs.Exists(path.Join(dir, configFilename)) {
			continue
		}
		res, _, err := s.AddIfNeeded(ctx, dir, path.Join(pth, dir), u)
		if err != nil {
			return nil, errors.Wrapf(err, "unable to load external module [%s]", dir)
		}
		ret = append(ret, res...)
	}

	return ret, nil
}
