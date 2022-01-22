// Content managed by Project Forge, see [projectforge.md] for details.
package theme

import (
	"path/filepath"

	"github.com/pkg/errors"
	"go.uber.org/zap"

	"github.com/kyleu/projectforge/app/lib/filesystem"
	"github.com/kyleu/projectforge/app/util"
)

const KeyNew = "new"

type Service struct {
	root   string
	files  filesystem.FileLoader
	cache  Themes
	logger *zap.SugaredLogger
}

func NewService(files filesystem.FileLoader, logger *zap.SugaredLogger) *Service {
	return &Service{root: "themes", files: files, logger: logger}
}

func (s *Service) All() Themes {
	s.loadIfNeeded()
	return s.cache
}

func (s *Service) Clear() {
	s.cache = nil
}

func (s *Service) Get(theme string) *Theme {
	for _, t := range s.All() {
		if t.Key == theme {
			return t
		}
	}
	return ThemeDefault
}

func (s *Service) Save(t *Theme) error {
	if t.Key == ThemeDefault.Key {
		return errors.New("can't overwrite default theme")
	}
	if t.Key == "" || t.Key == KeyNew {
		t.Key = util.RandomString(12)
	}
	b := util.ToJSONBytes(t, true)
	err := s.files.WriteFile(filepath.Join(s.root, t.Key+".json"), b, filesystem.DefaultMode, true)
	if err != nil {
		s.logger.Warnf("can't save theme [%s]: %+v", t.Key, err)
	}
	s.cache = s.cache.Replace(t)
	return nil
}

func (s *Service) loadIfNeeded() {
	if s.cache == nil {
		s.cache = Themes{ThemeDefault}
		for _, key := range s.files.ListJSON(s.root, true) {
			t := &Theme{}
			b, err := s.files.ReadFile(filepath.Join(s.root, key+".json"))
			if err != nil {
				s.logger.Warnf("can't load theme [%s]: %+v", key, err)
			}
			err = util.FromJSON(b, t)
			if err != nil {
				s.logger.Warnf("can't load theme [%s]: %+v", key, err)
			}
			t.Key = key
			s.cache = append(s.cache, t)
		}
		s.cache.Sort()
	}
}

func ApplyMap(frm util.ValueMap) *Theme {
	orig := ThemeDefault

	l := orig.Light.Clone().ApplyMap(frm, "light-")
	d := orig.Dark.Clone().ApplyMap(frm, "dark-")

	return &Theme{Light: l, Dark: d}
}
