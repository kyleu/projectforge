package module

import (
	"path/filepath"
	"sort"

	"github.com/kyleu/projectforge/app/filesystem"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

type Service struct {
	cache       map[string]*Module
	filesystems map[string]filesystem.FileLoader
	logger      *zap.SugaredLogger
}

func NewService(logger *zap.SugaredLogger) *Service {
	ret := &Service{cache: map[string]*Module{}, filesystems: map[string]filesystem.FileLoader{}, logger: logger}

	_, err := ret.LoadAll("core", "database", "desktop")
	if err != nil {
		logger.Errorf("unable to load [core] module: %+v", err)
	}
	return ret
}

func (s *Service) GetFilesystem(key string) filesystem.FileLoader {
	return filesystem.NewFileSystem(filepath.Join("module", key), s.logger)
}

func (s *Service) getNestedFilesystem(mods Modules) filesystem.FileLoader {
	var ret filesystem.FileLoader
	for _, mod := range mods {
		curr := s.GetFilesystem(mod.Key)
		if ret != nil {
			curr.AddChildren(ret)
		}
		ret = curr
	}
	return ret
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
	return ret, nil
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
