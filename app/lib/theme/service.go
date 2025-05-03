package theme

import (
	"path/filepath"

	"github.com/pkg/errors"
	"github.com/samber/lo"

	"projectforge.dev/projectforge/app/lib/filesystem"
	"projectforge.dev/projectforge/app/util"
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
	return lo.FindOrElse(s.All(logger), Default, func(t *Theme) bool {
		return t.Key == theme
	})
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
	err := s.files.WriteJSONFile(filepath.Join(s.root, t.Key+util.ExtJSON), t, filesystem.DefaultMode, true)
	if err != nil {
		logger.Warnf("can't save theme [%s]: %+v", t.Key, err)
	}
	s.cache = s.cache.Replace(t)
	return nil
}

func (s *Service) loadIfNeeded(logger util.Logger) {
	if s.cache == nil {
		s.cache = Themes{Default}
		if s.files.IsDir(s.root) {
			lo.ForEach(s.files.ListJSON(s.root, nil, true, logger), func(key string, _ int) {
				b, err := s.files.ReadFile(filepath.Join(s.root, key+util.ExtJSON))
				if err != nil {
					logger.Warnf("can't load theme [%s]: %+v", key, err)
				}
				t, err := util.FromJSONObj[*Theme](b)
				if err != nil {
					logger.Warnf("can't load theme [%s]: %+v", key, err)
				}
				t.Key = key
				s.cache = append(s.cache, t)
			})
		}
		s.cache = append(s.cache, lo.Reject(CatalogThemes, func(t *Theme, _ int) bool {
			return s.cache.Contains(t.Key)
		})...)
		s.cache.Sort()
	}
}

func (s *Service) Remove(key string, logger util.Logger) error {
	s.cache = s.cache.Remove(key)
	if !s.FileExists(key) {
		return nil
	}
	return s.files.Remove(filepath.Join(s.root, key+util.ExtJSON), logger)
}

func (s *Service) FileExists(key string) bool {
	return s.files.Exists(filepath.Join(s.root, key+util.ExtJSON))
}

func ApplyMap(frm util.ValueMap) *Theme {
	orig := Default

	l := orig.Light.Clone().ApplyMap(frm, "light-")
	d := orig.Dark.Clone().ApplyMap(frm, "dark-")

	return &Theme{Light: l, Dark: d}
}
