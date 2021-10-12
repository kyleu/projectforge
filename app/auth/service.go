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
