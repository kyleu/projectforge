// Content managed by Project Forge, see [projectforge.md] for details.
package cmd

import (
	"fmt"
	"runtime"

	"github.com/muesli/coral"
	"github.com/pkg/errors"
	"github.com/valyala/fasthttp"
	"go.uber.org/zap"

	"projectforge.dev/projectforge/app"
	"projectforge.dev/projectforge/app/controller"
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

	_, err = listenandserve(util.AppName, flags.Address, flags.Port, r)
	return err
}

func loadSite(flags *Flags, logger *zap.SugaredLogger) (fasthttp.RequestHandler, *zap.SugaredLogger, error) {
	r := controller.SiteRoutes()
	f := filesystem.NewFileSystem(flags.ConfigDir, logger)

	telemetryEnabled := util.GetEnv("disable_telemetry", "") == "true"
	st, err := app.NewState(flags.Debug, _buildInfo, f, telemetryEnabled, logger)
	if err != nil {
		return nil, logger, err
	}

	controller.SetSiteState(st)
	logger.Infof("started marketing site using address [%s:%d] on %s:%s", flags.Address, flags.Port, runtime.GOOS, runtime.GOARCH)
	return r, logger, nil
}
