package module

import (
	"path/filepath"
	"sort"

	"github.com/pkg/errors"
	"go.uber.org/zap"

	"github.com/kyleu/projectforge/app/export"
	"github.com/kyleu/projectforge/app/lib/filesystem"
	"github.com/kyleu/projectforge/app/lib/search/result"
	"github.com/kyleu/projectforge/app/project"
)

var nativeModuleKeys = []string{
	"android", "audit", "core", "database", "desktop", "export", "ios",
	"marketing", "migration", "mysql", "notarize", "oauth", "postgres",
	"sandbox", "schema", "search", "sqleditor", "sqlite", "types", "upgrade",
}

type Service struct {
	local       filesystem.FileLoader
	config      filesystem.FileLoader
	cache       map[string]*Module
	filesystems map[string]filesystem.FileLoader
	expSvc      *export.Service
	logger      *zap.SugaredLogger
}

func NewService(config filesystem.FileLoader, logger *zap.SugaredLogger) *Service {
	logger = logger.With("svc", "module")
	local := filesystem.NewFileSystem("module", logger)
	config = filesystem.NewFileSystem(filepath.Join(config.Root(), "module"), logger)
	fs := map[string]filesystem.FileLoader{}
	es := export.NewService(logger)
	ret := &Service{local: local, config: config, cache: map[string]*Module{}, filesystems: fs, expSvc: es, logger: logger}

	_, err := ret.LoadNative(nativeModuleKeys...)
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

func (s *Service) AddIfNeeded(key string, path string, url string) (*Module, bool, error) {
	ret, ok := s.cache[key]
	if ok {
		if ret.URL != url {
			s.logger.Warnf("module [%s] is loaded with url [%s] but there is another reference with url [%s]", ret.Key, ret.URL, url)
		}
		return ret, false, nil
	}
	m, err := s.load(key, path, url)
	if err != nil {
		return nil, false, err
	}
	s.cache[key] = m
	return m, true, nil
}

func (s *Service) Get(key string) (*Module, error) {
	ret, ok := s.cache[key]
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
	keys := make([]string, 0, len(s.cache))
	for k := range s.cache {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

func (s *Service) Modules() Modules {
	keys := s.Keys()
	ret := make(Modules, 0, len(keys))
	for _, key := range keys {
		p, _ := s.Get(key)
		ret = append(ret, p)
	}
	return ret
}

func (s *Service) Register(path string, defs ...*project.ModuleDef) ([]string, error) {
	var ret []string
	for _, def := range defs {
		_, added, err := s.AddIfNeeded(def.Key, filepath.Join(path, def.Path), def.URL)
		if err != nil {
			return nil, err
		}
		if added {
			ret = append(ret, def.Key)
		}
	}
	return ret, nil
}

func (s *Service) Search(q string) (result.Results, error) {
	ret := result.Results{}
	for _, mod := range s.Modules() {
		if res := result.NewResult("module", mod.Key, mod.WebPath(), mod.Title(), mod.IconSafe(), mod, q); len(res.Matches) > 0 {
			ret = append(ret, res)
		}
	}
	return ret, nil
}
