// Content managed by Project Forge, see [projectforge.md] for details.
package cmd

import (
	"projectforge.dev/projectforge/app"
	"projectforge.dev/projectforge/app/util"
)

func Run(bi *app.BuildInfo) (util.Logger, error) {
	_buildInfo = bi

	if err := rootCmd().Execute(); err != nil {
		return _logger, err
	}
	return _logger, nil
}
