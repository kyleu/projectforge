package auth

import (
	"github.com/samber/lo"

	"{{{ .Package }}}/app/util"
)

type Providers []*Provider

func (p Providers) Get(id string) *Provider {
	return lo.FindOrElse(p, nil, func(x *Provider) bool {
		return x.ID == id
	})
}

func (p Providers) Contains(id string) bool {
	return p.Get(id) != nil
}

func (p Providers) IDs() []string {
	return lo.Map(p, func(x *Provider, _ int) string {
		return x.ID
	})
}

func (p Providers) Titles() []string {
	return lo.Map(p, func(x *Provider, _ int) string {
		return x.Title
	})
}

func (s *Service) Providers(logger util.Logger) (Providers, error) {
	if s.providers == nil {
		err := s.load(logger)
		if err != nil {
			return nil, err
		}
	}
	return s.providers, nil
}
