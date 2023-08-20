package main

import (
	"os"

	"{{{ .Package }}}/app/cmd"
	"{{{ .Package }}}/app/lib/log"
)

func main() {
	logger, err := cmd.WASM()
	if err != nil {
		const msg = "exiting due to error"
		if logger == nil {
			println(log.Red.Add(err.Error())) //nolint:forbidigo
			println(log.Red.Add(msg))         //nolint:forbidigo
		} else {
			logger.Error(err)
			logger.Error(msg)
		}
		os.Exit(1)
	}
}
