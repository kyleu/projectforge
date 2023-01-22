package project

import (
	"os"
	"path/filepath"

	"github.com/pkg/errors"

	"projectforge.dev/projectforge/app/lib/filesystem"
	"projectforge.dev/projectforge/app/util"
)

func (s *Service) add(path string, parent *Project, logger util.Logger) (*Project, error) {
	if parent != nil {
		path = filepath.Join(parent.Path, path)
	}
	p, err := s.load(path, logger)
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
	s.cacheLock.Unlock()
	return p, nil
}

func (s *Service) getAdditional(logger util.Logger) ([]string, bool) {
	additionalContent, err := os.ReadFile(additionalFilename)
	if err != nil {
		return nil, false
	}
	var additional []string
	err = util.FromJSON(additionalContent, &additional)
	if err != nil {
		logger.Warnf("unable to parse additional projects from [%s]: %+v", additionalFilename, err)
	}
	return additional, true
}

func (s *Service) addAdditional(path string, logger util.Logger) {
	add, ok := s.getAdditional(logger)
	if !ok {
		return
	}
	hit := false
	for _, a := range add {
		if a == path {
			hit = true
			break
		}
	}
	if !hit {
		add = append(add, path)
		_ = os.WriteFile(additionalFilename, util.ToJSONBytes(add, true), filesystem.DefaultMode)
	}
}
