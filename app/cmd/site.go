// Package cmd - Content managed by Project Forge, see [projectforge.md] for details.
package cmd

import (
	"fmt"
	"net/http"
	"runtime"

	"github.com/muesli/coral"
	"github.com/pkg/errors"

	"projectforge.dev/projectforge/app"
	"projectforge.dev/projectforge/app/controller"
	"projectforge.dev/projectforge/app/controller/routes"
	"projectforge.dev/projectforge/app/lib/filesystem"
	"projectforge.dev/projectforge/app/util"
)

const keySite = "site"

func siteCmd() *coral.Command {
	short := fmt.Sprintf("Starts the marketing site on port %d (by default)", util.AppPort)
	f := func(*coral.Command, []string) error { return startSite(_flags) }
	ret := &coral.Command{Use: keySite, Short: short, RunE: f}
	return ret
}

func startSite(flags *Flags) error {
	err := initIfNeeded()
	if err != nil {
		return errors.Wrap(err, "error initializing application")
	}

	r, _, err := loadSite(flags, _logger)
	if err != nil {
		return err
	}

	_, err = listenandserve(flags.Address, flags.Port, r)
	return err
}

func loadSite(flags *Flags, logger util.Logger) (http.Handler, util.Logger, error) {
	r, err := routes.SiteRoutes(logger)
	if err != nil {
		return nil, logger, err
	}
	f, err := filesystem.NewFileSystem(flags.ConfigDir, false, "")
	if err != nil {
		return nil, logger, err
	}

	telemetryDisabled := util.GetEnvBool("disable_telemetry", false)
	st, err := app.NewState(flags.Debug, _buildInfo, f, !telemetryDisabled, flags.Port, logger)
	if err != nil {
		return nil, logger, err
	}

	controller.SetSiteState(st, logger)
	logger.Infof("started marketing site using address [%s:%d] on %s:%s", flags.Address, flags.Port, runtime.GOOS, runtime.GOARCH)
	return r, logger, nil
}
