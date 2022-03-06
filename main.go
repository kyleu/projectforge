// Content managed by Project Forge, see [projectforge.md] for details.
package main

import (
	"os"

	"projectforge.dev/app"
	"projectforge.dev/app/cmd"
	"projectforge.dev/app/lib/log"
)

var (
	version = "0.2.4" // updated by bin/tag.sh and ldflags
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
