// Content managed by Project Forge, see [projectforge.md] for details.
// Content managed by Project Forge, see [projectforge.md] for details.
package cmd

import (
	"projectforge.dev/projectforge/app"
	"projectforge.dev/projectforge/app/util"
)

func WASM() (util.Logger, error) {
	_buildInfo = &app.BuildInfo{Version: "-wasm", Commit: "wasm", Date: "unknown"}

	if err := rootCmd().Execute(); err != nil {
		return _logger, err
	}
	return _logger, nil
}
