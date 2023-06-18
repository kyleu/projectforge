package project

import (
	"github.com/samber/lo"
	"os"
	"path/filepath"
	"strings"

	"github.com/pkg/errors"

	"projectforge.dev/projectforge/app/lib/filesystem"
	"projectforge.dev/projectforge/app/util"
)

func (s *Service) add(path string, parent *Project) (*Project, error) {
	if parent != nil && !strings.HasPrefix(path, "/") {
		path = filepath.Join(parent.Path, path)
	}
	b, p, err := s.load(path)
	if err != nil {
		return nil, err
	}
	if parent != nil {
		p.Parent = parent.Key
	}
	s.cacheLock.Lock()
	curr, ok := s.cache[p.Key]
	s.cacheLock.Unlock()
	if ok {
		return nil, errors.Errorf("can't overwrite cache entry for project [%s] located at [%s]", curr.Key, curr.Path)
	}
	s.cacheLock.Lock()
	s.cache[p.Key] = p
	s.fileContent[p.Key] = b
	s.cacheLock.Unlock()
	return p, nil
}

func (s *Service) getAdditionalFilename() string {
	ret := ".projectforge/additional-projects.json"
	if _, err := os.Stat(ret); err == nil {
		return ret
	}
	return s.additional
}

func (s *Service) getAdditional(logger util.Logger) ([]string, bool) {
	additionalContent, err := os.ReadFile(s.getAdditionalFilename())
	if err != nil {
		return nil, false
	}
	var additional []string
	err = util.FromJSON(additionalContent, &additional)
	if err != nil {
		logger.Warnf("unable to parse additional projects from [%s]: %+v", s.getAdditionalFilename(), err)
	}
	return additional, true
}

func (s *Service) addAdditional(path string, logger util.Logger) {
	add, ok := s.getAdditional(logger)
	if !ok {
		return
	}
	hit := lo.ContainsBy(add, func(a string) bool {
		return a == path
	})
	if !hit {
		add = append(add, path)
		_ = os.WriteFile(s.getAdditionalFilename(), util.ToJSONBytes(add, true), filesystem.DefaultMode)
	}
}
