package auth

import (
	"go.uber.org/zap"
)

type Service struct {
	baseURL   string
	providers Providers
	logger    *zap.SugaredLogger
}

func NewService(baseURL string, logger *zap.SugaredLogger) *Service {
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
