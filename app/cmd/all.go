// Package cmd - Content managed by Project Forge, see [projectforge.md] for details.
package cmd

import (
	"fmt"

	"github.com/muesli/coral"
	"github.com/pkg/errors"

	"projectforge.dev/projectforge/app/util"
)

const keyAll = "all"

func allCmd() *coral.Command {
	short := fmt.Sprintf("Starts the main http server on port %d and the marketing site on port %d", util.AppPort, util.AppPort+1)
	f := func(*coral.Command, []string) error { return allF() }
	ret := &coral.Command{Use: keyAll, Short: short, RunE: f}
	return ret
}

func allF() error {
	if err := initIfNeeded(); err != nil {
		return errors.Wrap(err, "error initializing application")
	}

	go func() {
		if err := startSite(_flags.Clone(_flags.Port + 1)); err != nil {
			_logger.Errorf("unable to start marketing site: %+v", err)
		}
	}()
	return startServer(_flags)
}
