//go:build netbsd

package system

import (
	"context"
	"runtime"

	"{{{ .Package }}}/app/util"
)

func (s *Service) Status(ctx context.Context, logger util.Logger) (*Status, error) {
	logger.Warnf("calculating system status not supported on [%s]", runtime.GOOS)
	return &Status{}, nil
}
