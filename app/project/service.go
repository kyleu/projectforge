package project

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"sync"

	"github.com/kyleu/projectforge/app/filesystem"
	"github.com/kyleu/projectforge/app/util"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

const ConfigFilename = ".projectforge.json"

type Service struct {
	rootProject string
	cache       map[string]*Project
	cacheLock   sync.RWMutex
	filesystems map[string]filesystem.FileLoader
	fsLock   sync.RWMutex
	logger      *zap.SugaredLogger
}

func NewService(logger *zap.SugaredLogger) *Service {
	return &Service{cache: map[string]*Project{}, filesystems: map[string]filesystem.FileLoader{}, logger: logger}
}

func (s *Service) GetFilesystem(prj *Project) filesystem.FileLoader {
	curr, ok := s.filesystems[prj.Key]
	if ok {
		return curr
	}
	fs := filesystem.NewFileSystem(prj.Path, s.logger)
	s.fsLock.Lock()
	s.filesystems[prj.Key] = fs
	s.fsLock.Unlock()
	return fs
}

func (s *Service) add(path string, parent *Project) (*Project, error) {
	p, err := s.load(path)
	if err != nil {
		return nil, err
	}
	if parent != nil {
		p.Parent = parent.Key
	}
	curr, ok := s.cache[p.Key]
	if ok {
		return nil, errors.Errorf("can't overwrite cache entry for project [%s] located at [%s]", curr.Key, curr.Path)
	}
	s.cacheLock.Lock()
	s.cache[p.Key] = p
	s.cacheLock.Unlock()
	for _, kidKey := range p.Children {
		_, err := s.add(kidKey, p)
		if err != nil {
			return nil, errors.Wrapf(err, "error loading child [%s]", kidKey)
		}
	}
	return p, nil
}

func (s *Service) Refresh() (Projects, error) {
	s.cache = map[string]*Project{}
	p, err := s.add(".", nil)
	if err != nil {
		return nil, err
	}
	s.rootProject = p.Key
	return s.Projects(), nil
}

func (s *Service) Get(key string) (*Project, error) {
	ret, ok := s.cache[key]
	if ok {
		return ret, nil
	}
	return nil, errors.Errorf("no project with key [%s] found among %d candidates [%s]", key, len(s.cache), strings.Join(s.Keys(), ", "))
}

func (s *Service) Keys() []string {
	keys := make([]string, 0, len(s.cache))
	for k := range s.cache {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
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
		l, r := util.SplitStringLast(path, '/', true)
		if r == "" {
			r = l
		}
		ret := NewProject(r, path)
		ret.Name = r + " (missing)"
		return ret, nil
	}
	b, err := ioutil.ReadFile(cfgPath)
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
	tgtFS := s.GetFilesystem(prj)
	j := util.ToJSON(prj) + "\n"
	err := tgtFS.WriteFile(ConfigFilename, []byte(j), filesystem.DefaultMode, true)
	if err != nil {
		return errors.Wrapf(err, "unable to write config file to [%s]", ConfigFilename)
	}
	return nil
}
