package git

import (
	"projectforge.dev/projectforge/app/util"
)

const ok = "OK"

type Service struct {
	logger util.Logger
}

func NewService(logger util.Logger) *Service {
	logger = logger.With("svc", "build")
	return &Service{logger: logger}
}
