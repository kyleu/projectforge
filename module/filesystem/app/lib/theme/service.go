package theme

import (
	"path/filepath"

	"github.com/pkg/errors"
	"github.com/samber/lo"

	"{{{ .Package }}}/app/lib/filesystem"
	"{{{ .Package }}}/app/util"
)

const KeyNew = "new"

type Service struct {
	root  string
	files filesystem.FileLoader
	cache Themes
}

func NewService(files filesystem.FileLoader) *Service {
	return &Service{root: "themes", files: files}
}

func (s *Service) All(logger util.Logger) Themes {
	s.loadIfNeeded(logger)
	return s.cache
}

func (s *Service) Clear() {
	s.cache = nil
}

func (s *Service) Get(theme string, logger util.Logger) *Theme {
	return lo.FindOrElse(s.All(logger), ThemeDefault, func(t *Theme) bool {
		return t.Key == theme
	})
}

func (s *Service) Save(t *Theme, logger util.Logger) error {
	if t.Key == ThemeDefault.Key {
		return errors.New("can't overwrite default theme")
	}
	if t.Key == "" || t.Key == KeyNew {
		t.Key = util.RandomString(12)
	}
	b := util.ToJSONBytes(t, true)
	err := s.files.WriteFile(filepath.Join(s.root, t.Key+".json"), b, filesystem.DefaultMode, true)
	if err != nil {
		logger.Warnf("can't save theme [%s]: %+v", t.Key, err)
	}
	s.cache = s.cache.Replace(t)
	return nil
}

func (s *Service) loadIfNeeded(logger util.Logger) {
	if s.cache == nil {
		s.cache = Themes{ThemeDefault}{{{ if .HasModule "themecatalog" }}}
		s.cache = append(s.cache, CatalogThemes...){{{ end }}}
		lo.ForEach(s.files.ListJSON(s.root, nil, true, logger), func(key string, _ int) {
			t := &Theme{}
			b, err := s.files.ReadFile(filepath.Join(s.root, key+".json"))
			if err != nil {
				logger.Warnf("can't load theme [%s]: %+v", key, err)
			}
			err = util.FromJSON(b, t)
			if err != nil {
				logger.Warnf("can't load theme [%s]: %+v", key, err)
			}
			t.Key = key
			s.cache = append(s.cache, t)
		})
		s.cache.Sort()
	}
}

func ApplyMap(frm util.ValueMap) *Theme {
	orig := ThemeDefault

	l := orig.Light.Clone().ApplyMap(frm, "light-")
	d := orig.Dark.Clone().ApplyMap(frm, "dark-")

	return &Theme{Light: l, Dark: d}
}
