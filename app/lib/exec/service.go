package exec

import (
	"sync"
)

type Service struct {
	Execs Execs
	mu    sync.Mutex
}

func NewService() *Service {
	return &Service{}
}

func (s *Service) NewExec(key string, cmd string, path string, debug bool, envvars ...string) *Exec {
	s.mu.Lock()
	defer s.mu.Unlock()
	idx := len(s.Execs.GetByKey(key)) + 1
	e := NewExec(key, idx, cmd, path, debug, envvars...)
	s.Execs = append(s.Execs, e)
	s.Execs.Sort()
	return e
}
