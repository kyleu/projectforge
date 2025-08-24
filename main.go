package main // import projectforge.dev/projectforge

import (
	"context"

	"projectforge.dev/projectforge/app"
	"projectforge.dev/projectforge/app/cmd"
)

var (
	version = "1.7.20" // updated by bin/tag.sh and ldflags
	commit  = ""
	date    = "unknown"
)

func main() {
	cmd.Entrypoint(context.Background(), &app.BuildInfo{Version: version, Commit: commit, Date: date})
}
