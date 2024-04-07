package cmd

import (
	"os"

	"projectforge.dev/projectforge/app/lib/log"
)

func Entrypoint() {
	err := run()
	if err != nil {
		const msg = "exiting due to error"
		if _logger == nil {
			println(log.Red.Add(err.Error())) //nolint:forbidigo
			println(log.Red.Add(msg))         //nolint:forbidigo
		} else {
			_logger.Error(err)
			_logger.Error(msg)
		}
		os.Exit(1)
	}
}

func run() error {
	err := initIfNeeded()
	if err != nil {
		return err
	}
	return rootCmd().Execute()
}
