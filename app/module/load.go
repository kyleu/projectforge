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

func (s *Service) LoadNative(ctx context.Context, logger util.Logger, keys ...string) (Modules, error) {
	var ret Modules
	for _, key := range keys {
		m, err := s.Load(ctx, key, "", logger)
		if err != nil {
			return nil, err
		}
		ret = append(ret, m...)
	}
	return ret, nil
}

func (s *Service) Load(ctx context.Context, key string, url string, logger util.Logger) (Modules, error) {
	return s.load(ctx, key, "", url, logger)
}

func (s *Service) ConfigDirectory() string {
	return s.config.Root()
}

func (s *Service) load(ctx context.Context, key string, pth string, url string, logger util.Logger) (Modules, error) {
	var fs filesystem.FileLoader
	var err error
	switch {
	case s.local.IsDir(key):
		fs, err = filesystem.NewFileSystem(filepath.Join(s.local.Root(), key), false, "")
	case s.config.IsDir(key):
		fs, err = filesystem.NewFileSystem(filepath.Join(s.config.Root(), key), false, "")
	case pth != "":
		fs, err = filesystem.NewFileSystem(pth, false, "")
		if key == "*" {
			return s.loadDirectory(ctx, pth, url, fs, logger)
		}
	default:
		if url == "" {
			url, err = s.AssetURL(ctx, key, logger)
			if err != nil {
				return nil, errors.Wrapf(err, "error downloading module [%s]", key)
			}
		}
		err = s.Download(ctx, key, url, logger)
		if err != nil {
			return nil, errors.Wrapf(err, "error downloading module [%s]", key)
		}
		fs, err = filesystem.NewFileSystem(filepath.Join(s.config.Root(), key), false, "")
	}
	if err != nil {
		return nil, errors.Wrapf(err, "unable to load project file [%s]", key)
	}

	if !fs.Exists(configFilename) {
		const msg = "file [%s] does not exist in path for module [%s] using root [%s]"
		return nil, errors.Errorf(msg, configFilename, key, fs.Root())
	}

	b, err := fs.ReadFile(configFilename)
	if err != nil {
		return nil, errors.Wrapf(err, "unable to read [%s]", configFilename)
	}

	ret, err := util.FromJSONObj[*Module](b)
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

func (s *Service) loadDirectory(ctx context.Context, pth string, u string, fs filesystem.FileLoader, logger util.Logger) (Modules, error) {
	dirs := fs.ListDirectories(pth, nil, logger)
	if len(dirs) == 0 {
		return nil, errors.Errorf("directory at path [%s] does not contain module directories", pth)
	}
	ret := make(Modules, 0, len(dirs))
	for _, dir := range dirs {
		if !fs.Exists(path.Join(dir, configFilename)) {
			continue
		}
		res, _, err := s.AddIfNeeded(ctx, dir, path.Join(pth, dir), u, logger)
		if err != nil {
			return nil, errors.Wrapf(err, "unable to load external module [%s]", dir)
		}
		ret = append(ret, res...)
	}

	return ret, nil
}
