package main // import projectforge.dev/projectforge

import (
	"projectforge.dev/projectforge/app"
	"projectforge.dev/projectforge/app/cmd"
)

var (
	version = "1.7.7" // updated by bin/tag.sh and ldflags
	commit  = ""
	date    = "unknown"
)

func main() {
	cmd.Entrypoint(&app.BuildInfo{Version: version, Commit: commit, Date: date})
}
