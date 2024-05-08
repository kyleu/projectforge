package derive

import (
	"fmt"

	"projectforge.dev/projectforge/app/lib/metamodel/model"
	"projectforge.dev/projectforge/app/util"
)

type Section struct {
	Models model.Models `json:"models,omitempty"`
	Logs   []string     `json:"logs,omitempty"`
	Errors []string     `json:"errors,omitempty"`
}

func (s *Section) AddError(e error) *Section {
	s.Errors = append(s.Errors, e.Error())
	return s
}

func (s *Section) AddModel(m *model.Model) *Section {
	s.Models = append(s.Models, m)
	return s
}

func (s *Section) AddLog(msg string, args ...any) *Section {
	s.Logs = append(s.Logs, fmt.Sprintf(msg, args...))
	return s
}

type Result map[string]*Section

func (r Result) AddSection(key string, s *Section) Result {
	if _, ok := r[key]; ok {
		r[key+" [duplicate "+util.RandomString(4)+"]"] = s
	} else {
		r[key] = s
	}
	return r
}

func (r Result) AddErrorSection(key string, err error) Result {
	return r.AddSection(key, (&Section{}).AddError(err))
}
