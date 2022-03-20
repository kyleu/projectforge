package project

import (
	"os"
	"path/filepath"
	"sort"
	"strings"
	"sync"

	"github.com/pkg/errors"
	"go.uber.org/zap"
	"projectforge.dev/projectforge/app/lib/filesystem"
	"projectforge.dev/projectforge/app/util"
)

const (
	ConfigFilename     = ".projectforge.json"
	additionalFilename = "additional-projects.json"
)

type Service struct {
	cache       map[string]*Project
	cacheLock   sync.RWMutex
	filesystems map[string]filesystem.FileLoader
	fsLock      sync.RWMutex
	logger      *zap.SugaredLogger
}

func NewService(logger *zap.SugaredLogger) *Service {
	logger = logger.With("svc", "project")
	return &Service{cache: map[string]*Project{}, filesystems: map[string]filesystem.FileLoader{}, logger: logger}
}

func (s *Service) GetFilesystem(prj *Project) filesystem.FileLoader {
	s.fsLock.Lock()
	defer s.fsLock.Unlock()
	curr, ok := s.filesystems[prj.Key]
	if ok {
		return curr
	}
	fs := filesystem.NewFileSystem(prj.Path, s.logger)
	s.filesystems[prj.Key] = fs
	return fs
}

func (s *Service) add(path string, parent *Project) (*Project, error) {
	if parent != nil {
		path = filepath.Join(parent.Path, path)
	}
	p, err := s.load(path)
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

func (s *Service) Refresh() (Projects, error) {
	s.cache = map[string]*Project{}
	root, err := s.add(".", nil)
	if err != nil {
		return nil, err
	}
	if add, ok := s.getAdditional(); ok {
		for _, a := range add {
			if _, err := s.add(a, root); err != nil {
				return nil, err
			}
		}
	}
	return s.Projects(), nil
}

func (s *Service) Get(key string) (*Project, error) {
	s.cacheLock.Lock()
	ret, ok := s.cache[key]
	s.cacheLock.Unlock()
	if ok {
		return ret, nil
	}
	return nil, errors.Errorf("no project with key [%s] found among %d candidates [%s]", key, len(s.cache), strings.Join(s.Keys(), ", "))
}

func (s *Service) Keys() []string {
	s.cacheLock.Lock()
	keys := make(map[string]string, len(s.cache))
	titles := make([]string, 0, len(s.cache))
	for k, v := range s.cache {
		tl := strings.ToLower(v.Title())
		keys[tl] = k
		titles = append(titles, tl)
	}
	s.cacheLock.Unlock()
	sort.Strings(titles)
	ret := make([]string, 0, len(titles))
	for _, title := range titles {
		ret = append(ret, keys[title])
	}
	return ret
}

func (s *Service) Projects() Projects {
	keys := s.Keys()
	ret := make(Projects, 0, len(keys))
	for _, key := range keys {
		p, _ := s.Get(key)
		ret = append(ret, p)
	}
	return ret
}

func (s *Service) load(path string) (*Project, error) {
	cfgPath := filepath.Join(path, ConfigFilename)
	if curr, _ := os.Stat(cfgPath); curr == nil {
		l, r := util.StringSplitLast(path, '/', true)
		if r == "" {
			r = l
		}
		if r == "." {
			r, _ = os.Getwd()
			if strings.Contains(r, "/") {
				r = r[strings.LastIndex(r, "/")+1:]
			}
		}
		if r == "" {
			r = "root"
		}
		ret := NewProject(r, path)
		ret.Name = r + " (missing)"
		return ret, nil
	}
	b, err := os.ReadFile(cfgPath)
	if err != nil {
		return nil, err
	}

	ret := &Project{}
	err = util.FromJSON(b, &ret)
	if err != nil {
		return nil, errors.Wrapf(err, "can't load project from [%s]", ConfigFilename)
	}
	ret.Path = path

	return ret, nil
}

func (s *Service) Save(prj *Project) error {
	if prj.Icon == DefaultIcon {
		prj.Icon = ""
	}
	tgtFS := s.GetFilesystem(prj)
	j := util.ToJSON(prj) + "\n"
	err := tgtFS.WriteFile(ConfigFilename, []byte(j), filesystem.DefaultMode, true)
	if err != nil {
		return errors.Wrapf(err, "unable to write config file to [%s]", ConfigFilename)
	}
	if prj.Path != "" && prj.Path != "." {
		s.addAdditional(prj.Path)
	}
	return nil
}

func (s *Service) ByPath(path string) *Project {
	s.cacheLock.Lock()
	defer s.cacheLock.Unlock()
	for _, v := range s.cache {
		if v.Path == path {
			return v
		}
	}
	return nil
}

func (s *Service) Init() error {
	return nil
}

func (s *Service) getAdditional() ([]string, bool) {
	_, fs := s.root()
	additionalContent, err := fs.ReadFile(additionalFilename)
	if err != nil {
		s.logger.Warnf("unable to load additional projects from [%s]", filepath.Join(fs.Root(), additionalFilename))
	}
	var additional []string
	err = util.FromJSON(additionalContent, &additional)
	if err != nil {
		s.logger.Warnf("unable to parse additional projects from [%s]: %+v", filepath.Join(fs.Root(), additionalFilename), err)
	}
	return additional, true
}

func (s *Service) addAdditional(path string) {
	add, ok := s.getAdditional()
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
