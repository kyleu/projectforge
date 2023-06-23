package module

import (
	"context"
	"path/filepath"
	"sync"

	"github.com/pkg/errors"
	"github.com/samber/lo"
	"golang.org/x/exp/slices"

	"projectforge.dev/projectforge/app/lib/filesystem"
	"projectforge.dev/projectforge/app/lib/search/result"
	"projectforge.dev/projectforge/app/project/export"
	"projectforge.dev/projectforge/app/util"
)

var nativeModuleKeys = []string{
	"android", "audit", "core", "database", "databaseui", "desktop", "docbrowse", "export", "expression",
	"filesystem", "graphql", "ios", "jsx", "marketing", "migration", "mysql", "notarize", "oauth", "postgres",
	"process", "readonlydb", "sandbox", "schema", "search",
	"sqlite", "sqlserver", "types", "upgrade", "user", "wasm", "websocket",
}

type Service struct {
	local       filesystem.FileLoader
	config      filesystem.FileLoader
	cache       map[string]*Module
	cacheMu     sync.Mutex
	filesystems map[string]filesystem.FileLoader
	expSvc      *export.Service
}

func NewService(ctx context.Context, config filesystem.FileLoader, logger util.Logger) *Service {
	local := filesystem.NewFileSystem("module")
	config = filesystem.NewFileSystem(filepath.Join(config.Root(), "module"))
	fs := map[string]filesystem.FileLoader{}
	es := export.NewService()
	ret := &Service{local: local, config: config, cache: map[string]*Module{}, filesystems: fs, expSvc: es}

	_, err := ret.LoadNative(ctx, logger, nativeModuleKeys...)
	if err != nil {
		logger.Errorf("unable to load [core] module: %+v", err)
	}
	return ret
}

func (s *Service) GetFilesystem(key string) filesystem.FileLoader {
	mod, err := s.Get(key)
	if err != nil {
		if s.local.IsDir(key) {
			return s.local
		}
		return s.config
	}
	return mod.Files
}

func (s *Service) AddIfNeeded(ctx context.Context, key string, path string, url string, logger util.Logger) (Modules, bool, error) {
	s.cacheMu.Lock()
	ret, ok := s.cache[key]
	s.cacheMu.Unlock()
	if ok {
		if ret.URL != url {
			logger.Warnf("module [%s] is loaded with url [%s] but there is another reference with url [%s]", ret.Key, ret.URL, url)
		}
		return Modules{ret}, false, nil
	}
	mods, err := s.load(ctx, key, path, url, logger)
	if err != nil {
		return nil, false, err
	}
	lo.ForEach(mods, func(mod *Module, _ int) {
		s.cacheMu.Lock()
		s.cache[mod.Key] = mod
		s.cacheMu.Unlock()
	})
	return mods, true, nil
}

func (s *Service) Get(key string) (*Module, error) {
	s.cacheMu.Lock()
	ret, ok := s.cache[key]
	s.cacheMu.Unlock()
	if !ok {
		return nil, errors.New("no module available with key [" + key + "]")
	}
	return ret, nil
}

func (s *Service) GetModules(keys ...string) (Modules, error) {
	var ret Modules
	for _, m := range keys {
		mod, err := s.Get(m)
		if err != nil {
			return nil, err
		}
		ret = append(ret, mod)
	}
	return ret.Sort(), nil
}

func (s *Service) Keys() []string {
	s.cacheMu.Lock()
	defer s.cacheMu.Unlock()
	keys := lo.Keys(s.cache)
	slices.Sort(keys)
	return keys
}

func (s *Service) Modules() Modules {
	return lo.Map(s.Keys(), func(key string, _ int) *Module {
		p, _ := s.Get(key)
		return p
	})
}

func (s *Service) Deps() map[string][]string {
	return lo.Associate(s.Modules(), func(m *Module) (string, []string) {
		return m.Key, m.Requires
	})
}

func (s *Service) Register(ctx context.Context, root string, key string, path string, u string, logger util.Logger) ([]string, error) {
	var ret []string
	_, added, err := s.AddIfNeeded(ctx, key, filepath.Join(root, path), u, logger)
	if err != nil {
		return nil, err
	}
	if added {
		ret = append(ret, key)
	}
	return ret, nil
}

func (s *Service) Search(ctx context.Context, q string, logger util.Logger) (result.Results, error) {
	return lo.FilterMap(s.Modules(), func(mod *Module, _ int) (*result.Result, bool) {
		res := result.NewResult("module", mod.Key, mod.WebPath(), mod.Title(), mod.IconSafe(), mod, mod, q)
		if len(res.Matches) > 0 {
			return res, true
		}
		return nil, false
	}), nil
}
