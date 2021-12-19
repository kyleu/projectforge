package git

import (
	"fmt"

	"github.com/kyleu/projectforge/app/project"
)

type Status struct {
	Key     string           `json:"key"`
	Project *project.Project `json:"-"`
	Branch  string           `json:"branch"`
	Dirty   []string         `json:"dirty"`
	Error   string           `json:"error"`
}

func (s *Status) Status() string {
	if s.Dirty == nil {
		return "No repo"
	}
	if len(s.Dirty) > 0 {
		return fmt.Sprintf("[%d] changes", len(s.Dirty))
	}
	return "OK"
}

func (s *Status) Actions() Actions {
	ret := Actions{ActionStatus, ActionMagic}
	if s.Dirty == nil {
		return append(ret, ActionCreateRepo)
	}
	if len(s.Dirty) > 0 {
		ret = append(ret, ActionCommit)
	}
	return ret
}

type Statuses []*Status

func (s Statuses) Get(key string) *Status {
	for _, x := range s {
		if x.Project.Key == key {
			return x
		}
	}
	return nil
}
