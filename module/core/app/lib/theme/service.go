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
	return ThemeDefault
}

func (s *Service) Save(t *Theme, logger util.Logger) error {
	if t.Key == ThemeDefault.Key {
		return errors.New("can't overwrite default theme")
	}
	if t.Key == "" || t.Key == KeyNew {
		t.Key = util.RandomString(12)
	}
	s.cache = s.cache.Replace(t)
	return nil
}

func (s *Service) loadIfNeeded(logger util.Logger) {
	if s.cache == nil {
		s.cache = Themes{ThemeDefault}
	}
}
