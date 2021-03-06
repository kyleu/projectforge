package git

import (
	"strconv"

	"projectforge.dev/projectforge/app/project"
	"projectforge.dev/projectforge/app/util"
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
	ret := Actions{ActionStatus}
	if s.Status == "no repo" {
		return append(ret, ActionCreateRepo)
	}
	ret = append(ret, ActionFetch)
	if dirty := s.DataStringArray("dirty"); len(dirty) > 0 {
		ret = append(ret, ActionCommit)
	}
	if s.DataInt("commitsAhead") > 0 {
		ret = append(ret, ActionPush)
	}
	if s.DataInt("commitsBehind") > 0 {
		ret = append(ret, ActionPull)
	}
	ret = append(ret, ActionUndoCommit, ActionMagic)
	return ret
}

func (s *Result) DataString(k string) string {
	if s.Data == nil {
		return ""
	}
	return s.Data.GetStringOpt(k)
}

func (s *Result) DataInt(k string) int {
	ret, _ := strconv.Atoi(s.DataString(k))
	return ret
}

func (s *Result) DataStringArray(k string) []string {
	if s.Data == nil {
		return nil
	}
	ret, _ := s.Data.GetStringArray(k, true)
	return ret
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
