package git

import (
	"github.com/kyleu/projectforge/app/project"
	"github.com/kyleu/projectforge/app/util"
)

type Result struct {
	Project *project.Project `json:"-"`
	Status  string           `json:"status"`
	Data    util.ValueMap    `json:"data"`
	Error   string           `json:"error"`
}

func NewResult(prj *project.Project, status string, data util.ValueMap) *Result {
	return &Result{Project: prj, Status: status, Data: data}
}

func (s *Result) Actions() Actions {
	ret := Actions{ActionStatus, ActionMagic}
	if s.Status == "no repo" {
		return append(ret, ActionCreateRepo)
	}
	ret = append(ret, ActionFetch)
	if dirty, _ := s.Data.GetStringArray("dirty", true); len(dirty) > 0 {
		ret = append(ret, ActionCommit)
	}
	return ret
}

func (s *Result) Branch() string {
	if s.Data == nil {
		return ""
	}
	return s.Data.GetStringOpt("branch")
}

type Results []*Result

func (s Results) Get(key string) *Result {
	for _, x := range s {
		if x.Project.Key == key {
			return x
		}
	}
	return nil
}
