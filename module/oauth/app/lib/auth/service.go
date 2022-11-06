package auth

import (
	"{{{ .Package }}}/app/util"
)

type Service struct {
	baseURL   string
	port      uint16
	providers Providers
}

func NewService(baseURL string, port uint16, logger util.Logger) *Service {
	ret := &Service{baseURL: baseURL, port: port}
	_ = ret.load(logger)
	return ret
}

func (s *Service) LoginURL() string {
	if len(s.providers) == 1 {
		return "/auth/" + s.providers[0].ID
	}
	return defaultProfilePath
}
