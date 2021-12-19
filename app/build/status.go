package build

import (
	"github.com/kyleu/projectforge/app/project"
)

type Status struct {
	Key     string           `json:"key"`
	Project *project.Project `json:"-"`
	Git     interface{}      `json:"git"`
	Error   string           `json:"error"`
}

func (s *Status) Status() string {
	return "OK"
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
