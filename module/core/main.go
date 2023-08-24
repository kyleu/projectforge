package main // import {{{ .Package }}}

import (
	"{{{ .Package }}}/app"
	"{{{ .Package }}}/app/cmd"
)

var (
	version = "{{{ .Version }}}" // updated by bin/tag.sh and ldflags
	commit  = ""
	date    = "unknown"
)

func main() {
	cmd.Entrypoint(&app.BuildInfo{Version: version, Commit: commit, Date: date})
}
