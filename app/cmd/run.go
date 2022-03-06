// Content managed by Project Forge, see [projectforge.md] for details.
package cmd

import (
	"go.uber.org/zap"

	"projectforge.dev/app"
)

func Run(bi *app.BuildInfo) (*zap.SugaredLogger, error) {
	_buildInfo = bi

	if err := rootCmd().Execute(); err != nil {
		return _logger, err
	}
	return _logger, nil
}
