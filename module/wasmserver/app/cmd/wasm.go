package cmd

import (
	"{{{ .Package }}}/app"
	"{{{ .Package }}}/app/util"
)

func WASM() (util.Logger, error) {
	_buildInfo = &app.BuildInfo{Version: "-wasm", Commit: "wasm", Date: "unknown"}

	if err := rootCmd().Execute(); err != nil {
		return _logger, err
	}
	return _logger, nil
}
