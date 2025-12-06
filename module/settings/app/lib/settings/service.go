package settings

import (
	"sync"

	"github.com/pkg/errors"

	"{{{ .Package }}}/app/lib/filesystem"
	"{{{ .Package }}}/app/util"
)

const settingsPath = "settings.json"

type Service struct {
	cached *Settings
	mu     sync.Mutex
	fs     filesystem.FileLoader
}

func NewService(fs filesystem.FileLoader) *Service {
	return &Service{fs: fs}
}

func (s *Service) Sync() *Settings {
	s.cached = s.load()
	return s.cached.Clone()
}

func (s *Service) Get() *Settings {
	if s.cached == nil {
		return s.Sync()
	}
	return s.cached.Clone()
}

func (s *Service) Set(x *Settings) error {
	return s.save(x)
}

func (s *Service) SetMap(m util.ValueMap) error {
	ret, extra, err := SettingsFromMap(m, true)
	if err != nil {
		return err
	}
	if len(extra) > 0 {
		return errors.Errorf("unknown settings [%s] included in request", extra.JSON())
	}
	return s.Set(ret)
}

func (s *Service) load() *Settings {
	s.mu.Lock()
	defer s.mu.Unlock()
	ret := &Settings{}
	if s.fs.Exists(settingsPath) {
		b, err := s.fs.ReadFile(settingsPath)
		if err != nil {
			return ret
		}
		err = util.FromJSON(b, ret)
		if err != nil {
			return ret
		}
	}
	return ret
}

func (s *Service) save(x *Settings) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.cached = x.Clone()
	b := util.ToJSONBytes(x, true)
	return s.fs.WriteFile(settingsPath, b, filesystem.DefaultMode, true)
}
