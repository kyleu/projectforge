package auth

import (
	"{{{ .Package }}}/app/util"
)

type Service struct {
	baseURL   string
	providers Providers
	logger    util.Logger
}

func NewService(baseURL string, logger util.Logger) *Service {
	ret := &Service{baseURL: baseURL, logger: logger}
	_ = ret.load()
	return ret
}

func (s *Service) LoginURL() string {
	if len(s.providers) == 1 {
		return "/auth/" + s.providers[0].ID
	}
	return defaultProfilePath
}
