package service

import (
	"strings"

	"github.com/samber/lo"

	"projectforge.dev/projectforge/app/util"
)

const defaultIcon = "service"

type Service struct {
	Name        string                `json:"name"`
	Package     string                `json:"package"`
	Group       []string              `json:"group,omitempty"`
	Schema      string                `json:"schema,omitempty"`
	Description string                `json:"description,omitempty"`
	Icon        string                `json:"icon,omitempty"`
	Calls       *util.OrderedMap[any] `json:"calls,omitempty"`
	Config      any                   `json:"config,omitempty"`
}

func (s *Service) FirstLetter() any {
	return strings.ToLower(s.Name[0:1])
}

func (s *Service) IconSafe() string {
	if _, ok := util.SVGLibrary[s.Icon]; ok {
		return s.Icon
	}
	return defaultIcon
}

func (s *Service) Camel() string {
	return util.StringToLowerCamel(s.Name)
}

func (s *Service) CamelLower() string {
	return strings.ToLower(s.Camel())
}

type Services []*Service

func (s Services) Get(key string) *Service {
	return lo.FindOrElse(s, nil, func(x *Service) bool {
		return x.Name == key
	})
}
