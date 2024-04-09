package theme

import (
	"github.com/pkg/errors"

	"{{{ .Package }}}/app/util"
)

const KeyNew = "new"

type Service struct {
	root  string
	cache Themes
}

func NewService() *Service {
	return &Service{root: "themes"}
}

func (s *Service) All(logger util.Logger) Themes {
	s.loadIfNeeded(logger)
	return s.cache
}

func (s *Service) Clear() {
	s.cache = nil
}

func (s *Service) Get(theme string, logger util.Logger) *Theme {
	for _, t := range s.All(logger) {
		if t.Key == theme {
			return t
		}
	}
	return Default
}

func (s *Service) Save(t *Theme, originalKey string, logger util.Logger) error {
	if t.Matches(Default) {
		return errors.New("can't overwrite default theme")
	}
	if t.Key == "" || t.Key == KeyNew {
		t.Key = util.RandomString(12)
	}
	if originalKey != t.Key {
		err := s.Remove(originalKey, logger)
		if err != nil {
			return err
		}
	}
	s.cache = s.cache.Replace(t)
	return nil
}

func (s *Service) Remove(key string, logger util.Logger) error {
	s.cache = s.cache.Remove(key)
	return nil
}

func (s *Service) FileExists(key string) bool {
	return false
}

func (s *Service) loadIfNeeded(logger util.Logger) {
	if s.cache == nil {
		s.cache = Themes{Default}{{{ if .HasModule "themecatalog" }}}
		s.cache = append(s.cache, CatalogThemes...){{{ end }}}
	}
}
