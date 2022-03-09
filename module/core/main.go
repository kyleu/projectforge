package main // import {{{ .Package }}}

import (
	"os"

	"{{{ .Package }}}/app"
	"{{{ .Package }}}/app/cmd"
	"{{{ .Package }}}/app/lib/log"
)

var (
	version = "{{{ .Version }}}" // updated by bin/tag.sh and ldflags
	commit  = ""
	date    = "unknown"
)

func main() {
	logger, err := cmd.Run(&app.BuildInfo{Version: version, Commit: commit, Date: date})
	if err != nil {
		const msg = "exiting due to error"
		if logger == nil {
			println(log.Red.Add(err.Error())) //nolint
			println(log.Red.Add(msg))         //nolint
		} else {
			logger.Error(err)
			logger.Error(msg)
		}
		os.Exit(1)
	}
}
