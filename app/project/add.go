package project

import (
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
	_, fs := s.root()
	if fs == nil {
		return nil, false
	}
	additionalContent, err := fs.ReadFile(additionalFilename)
	if err != nil {
		logger.Warnf("unable to load additional projects from [%s]", filepath.Join(fs.Root(), additionalFilename))
		return nil, false
	}
	var additional []string
	err = util.FromJSON(additionalContent, &additional)
	if err != nil {
		logger.Warnf("unable to parse additional projects from [%s]: %+v", filepath.Join(fs.Root(), additionalFilename), err)
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
		_, fs := s.root()
		_ = fs.WriteFile(additionalFilename, util.ToJSONBytes(add, true), filesystem.DefaultMode, true)
	}
}

func (s *Service) root() (*Project, filesystem.FileLoader) {
	root := s.ByPath(".")
	if root == nil {
		return nil, nil
	}
	fs := s.GetFilesystem(root)
	if !fs.Exists(additionalFilename) {
		return nil, nil
	}
	return root, fs
}
