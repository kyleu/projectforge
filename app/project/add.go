package project

import (
	"path/filepath"
	"strings"

	"github.com/samber/lo"

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
	s.cache[p.Key] = p
	s.fileContent[p.Key] = b
	s.cacheLock.Unlock()
	return p, nil
}

func (s *Service) getAdditionalFilename(fs filesystem.FileLoader) string {
	ret := ".projectforge/additional-projects.json"
	if _, err := fs.Stat(ret); err == nil {
		return ret
	}
	return s.additional
}

func (s *Service) getAdditional(fs filesystem.FileLoader, createIfMissing bool, logger util.Logger) ([]string, bool) {
	pth := s.getAdditionalFilename(fs)
	additionalContent, err := fs.ReadFile(pth)
	if err != nil {
		if !createIfMissing {
			return nil, false
		}
		additionalContent = []byte("[]")
		err = fs.WriteFile(pth, additionalContent, filesystem.DefaultMode, false)
		if err != nil {
			return nil, false
		}
	}
	var additional []string
	err = util.FromJSON(additionalContent, &additional)
	if err != nil {
		logger.Warnf("unable to parse additional projects from [%s]: %+v", s.getAdditionalFilename(fs), err)
	}
	return additional, true
}

func (s *Service) addAdditional(path string, fs filesystem.FileLoader, logger util.Logger) {
	add, ok := s.getAdditional(fs, true, logger)
	if !ok {
		return
	}
	hit := lo.ContainsBy(add, func(a string) bool {
		return a == path
	})
	if !hit {
		add = append(add, path)
		_ = fs.WriteFile(s.getAdditionalFilename(fs), util.ToJSONBytes(add, true), filesystem.DefaultMode, true)
	}
}
