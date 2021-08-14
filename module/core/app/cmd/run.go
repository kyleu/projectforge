package cmd

import (
	"go.uber.org/zap"

	"{{{ .Package }}}/app"
)

func Run(bi *app.BuildInfo) (*zap.SugaredLogger, error) {
	_buildInfo = bi

	if err := rootCmd().Execute(); err != nil {
		return _logger, err
	}
	return _logger, nil
}